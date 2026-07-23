package admin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"app/pkg/agent"
	"app/pkg/doubaoasr"
	"app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/server/response"
	"app/service/mcptool"
	"app/service/motto"
	"app/store"

	"github.com/openai/openai-go/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.AIServiceHTTPServer = (*AI)(nil)

const maxRemindSpeechAudioBytes = 20 * 1024 * 1024

type remindSpeechTranscriber interface {
	Transcribe(ctx context.Context, audioBase64 string) (string, error)
}

type AI struct {
	agent             *agent.Agent
	store             *store.Store
	asrConf           doubaoasr.Config
	speechTranscriber remindSpeechTranscriber
}

func NewAI(aiAgent *agent.Agent, s *store.Store, asrConf doubaoasr.Config) *AI {
	return &AI{
		agent:   aiAgent,
		store:   s,
		asrConf: asrConf,
	}
}

// ChatMessage 表示一次对话中的单条消息。
type ChatMessage struct {
	Role            string            `json:"role"`
	Content         string            `json:"content"`
	ContextMessages []json.RawMessage `json:"contextMessages,omitempty"`
}

// ChatRequest 表示携带历史消息的聊天请求。
type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// ToolStartEvent 表示发送给前端的工具调用开始事件。
type ToolStartEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MCPName   string `json:"mcpName"`
	Arguments string `json:"arguments"`
}

// ToolEndEvent 表示发送给前端的工具调用结束事件。
type ToolEndEvent struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

// StreamEvent 表示发送给前端的 SSE 流事件
type StreamEvent struct {
	Type    string `json:"type"`              // content, reasoning, tool_start, tool_end, context, error, done
	Content string `json:"content,omitempty"` // 用于 text content 和 error message
	Data    any    `json:"data,omitempty"`    // 用于其他事件类型的数据载荷
}

func (a *AI) sendStreamEvent(w http.ResponseWriter, flusher http.Flusher, event StreamEvent) {
	b, _ := json.Marshal(event)
	fmt.Fprintf(w, "data: %s\n\n", string(b))
	flusher.Flush()
}

// Chat 处理 SSE 流式 AI 聊天响应。
func (a *AI) Chat(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errors.BadRequest("INVALID_REQUEST", "Invalid request body"))
		return
	}

	if len(req.Messages) == 0 {
		response.Fail(w, errors.BadRequest("EMPTY_MESSAGES", "Messages cannot be empty"))
		return
	}

	// Validate that there's at least one user message with content
	hasContent := false
	for _, msg := range req.Messages {
		if strings.TrimSpace(msg.Content) != "" {
			hasContent = true
			break
		}
	}
	if !hasContent {
		response.Fail(w, errors.BadRequest("EMPTY_MESSAGE", "Message content cannot be empty"))
		return
	}

	// Set SSE headers early
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		response.Fail(w, errors.InternalServer("STREAMING_ERROR", "Streaming not supported"))
		return
	}

	ctx := r.Context()

	// Build system prompt
	prompt := a.buildSystemPrompt()

	// 根据请求历史构造 OpenAI 消息。
	openAIMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, msg := range req.Messages {
		if len(msg.ContextMessages) > 0 {
			for _, rawMessage := range msg.ContextMessages {
				openAIMessage, err := agent.DecodeMessageParam(rawMessage)
				if err != nil {
					response.Fail(w, errors.BadRequest("INVALID_CONTEXT_MESSAGE", "Invalid context message"))
					return
				}
				openAIMessages = append(openAIMessages, openAIMessage)
			}
			continue
		}

		if strings.TrimSpace(msg.Content) == "" {
			continue
		}
		switch msg.Role {
		case "user":
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Content))
		case "assistant":
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Content))
		}
	}

	result, err := a.agent.Run(ctx, agent.Request{
		SystemPrompt: prompt,
		Messages:     openAIMessages,
		UseTools:     true,
	}, agent.EventHandler{
		OnContent: func(_ context.Context, content string) error {
			a.sendStreamEvent(w, flusher, StreamEvent{Type: "content", Content: content})
			return nil
		},
		OnReasoning: func(_ context.Context, content string) error {
			a.sendStreamEvent(w, flusher, StreamEvent{Type: "reasoning", Content: content})
			return nil
		},
		OnToolStart: func(_ context.Context, event agent.ToolEvent) error {
			a.sendStreamEvent(w, flusher, StreamEvent{
				Type: "tool_start",
				Data: ToolStartEvent{
					ID:        event.ID,
					Name:      event.Name,
					MCPName:   event.MCPName,
					Arguments: event.Arguments,
				},
			})
			return nil
		},
		OnToolEnd: func(_ context.Context, event agent.ToolEvent) error {
			a.sendStreamEvent(w, flusher, StreamEvent{
				Type: "tool_end",
				Data: ToolEndEvent{
					ID:     event.ID,
					Result: event.Result,
				},
			})
			return nil
		},
	})
	if err != nil {
		if ctx.Err() == context.Canceled {
			return
		}
		a.sendStreamEvent(w, flusher, StreamEvent{Type: "error", Content: err.Error()})
		return
	}

	a.sendStreamEvent(w, flusher, StreamEvent{Type: "context", Data: result.Messages})

	// Send done event
	a.sendStreamEvent(w, flusher, StreamEvent{Type: "done"})
}

// buildSystemPrompt 构造带工具调用说明的系统提示词。
func (a *AI) buildSystemPrompt() string {
	basePrompt := fmt.Sprintf(`You are a helpful assistant. Respond in the same language as the user's message.
IMPORTANT: Any deep thinking, reasoning, or thought processes MUST be output in the same language as the user's message.

When you encounter questions that you cannot answer directly, such as:
- Current events, news, or real-time information
- Recent research or developments in professional fields
- Specific facts you are uncertain about
- Information that may have changed after your training data cutoff

You should use the available tools to find accurate and up-to-date information.

Current Time: %s
`, time.Now().Format(time.DateTime))

	return basePrompt
}

func (a *AI) GenerateTags(ctx context.Context, req *adminv1.GenerateTagsRequest) (*adminv1.GenerateTagsResponse, error) {
	content := strings.TrimSpace(req.GetContent())
	if content == "" {
		return nil, errors.BadRequest("EMPTY_CONTENT", "Content cannot be empty")
	}

	aiClient := a.agent.GetClient(ctx)
	aiModel := a.agent.GetModel(ctx)

	prompt := `你是一个博客写作助手。请根据用户提供的文章标题与正文，为文章生成 3-8 个中文标签。
要求：
1) 标签要简短（2-6 个字），可用中英文混合（如 Go、React、MySQL）。
2) 标签去重，不要包含无意义的词（比如"文章""随笔""记录"）。
3) 只输出 JSON 数组（例如：["Go","数据库","性能优化"]），不要输出其它任何文字。`

	userInput := fmt.Sprintf("标题：%s\n正文：\n%s", strings.TrimSpace(req.GetTitle()), content)

	aiReq := openai.ChatCompletionNewParams{
		Model: openai.ChatModel(aiModel),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userInput),
		},
	}

	completion, err := aiClient.Chat.Completions.New(ctx, aiReq)
	if err != nil {
		return nil, errors.InternalServer("AI_GENERATE_TAGS_ERROR", err.Error())
	}
	if len(completion.Choices) == 0 {
		return adminv1.GenerateTagsResponse_builder{Tags: []string{}}.Build(), nil
	}

	raw := strings.TrimSpace(completion.Choices[0].Message.Content)
	tags := parseTagsFromAIResponse(raw)
	return adminv1.GenerateTagsResponse_builder{Tags: tags}.Build(), nil
}

func (a *AI) RemindSmartCreate(ctx context.Context, req *adminv1.RemindSmartCreateRequest) (*types.IDResponse, error) {
	aiClient := a.agent.GetClient(ctx)
	aiModel := a.agent.GetModel(ctx)

	lastID, err := mcptool.SmartCreateRemind(ctx, aiClient, aiModel, a.store, req.GetContent())
	if err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(lastID)}.Build(), nil
}

func (a *AI) RemindSpeechTranscribe(ctx context.Context, req *adminv1.RemindSpeechTranscribeRequest) (*adminv1.RemindSpeechTranscribeResponse, error) {
	audioBase64 := strings.TrimSpace(req.GetAudioBase64())
	if audioBase64 == "" {
		return nil, errors.BadRequest("EMPTY_AUDIO", "录音不能为空")
	}

	audioData, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return nil, errors.BadRequest("INVALID_AUDIO_BASE64", "录音数据格式错误").WithCause(err)
	}
	if len(audioData) == 0 {
		return nil, errors.BadRequest("EMPTY_AUDIO", "录音不能为空")
	}
	if len(audioData) > maxRemindSpeechAudioBytes {
		return nil, errors.BadRequest("AUDIO_TOO_LARGE", "录音文件过大，请缩短录音时间")
	}

	text, err := a.getSpeechTranscriber().Transcribe(ctx, audioBase64)
	if err != nil {
		return nil, errors.InternalServer("AI_REMIND_TRANSCRIBE_ERROR", "语音识别失败").WithCause(err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.BadRequest("EMPTY_TRANSCRIPT", "没有识别到文字，请重新录音")
	}
	return adminv1.RemindSpeechTranscribeResponse_builder{Text: text}.Build(), nil
}

func (a *AI) GenerateMood(ctx context.Context, _ *emptypb.Empty) (*adminv1.GenerateMoodResponse, error) {
	dateStr := time.Now().Format("2006-01-02")
	provider := motto.NewOpenAIProvider(a.agent)
	content, err := provider.Generate(ctx, dateStr)
	if err != nil {
		return nil, err
	}

	return adminv1.GenerateMoodResponse_builder{Content: content}.Build(),
		nil
}

func (a *AI) getSpeechTranscriber() remindSpeechTranscriber {
	if a.speechTranscriber != nil {
		return a.speechTranscriber
	}

	return doubaoasr.Client{
		APIKey:     firstNonEmpty(os.Getenv("DOUBAO_ASR_API_KEY"), a.asrConf.APIKey),
		Endpoint:   firstNonEmpty(os.Getenv("DOUBAO_ASR_ENDPOINT"), a.asrConf.Endpoint),
		ResourceID: firstNonEmpty(os.Getenv("DOUBAO_ASR_RESOURCE_ID"), a.asrConf.ResourceID),
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func parseTagsFromAIResponse(text string) []string {
	if tags, ok := parseTagsFromJSONArray(text); ok {
		return normalizeTags(tags)
	}

	extracted, ok := extractFirstJSONArray(text)
	if !ok {
		return []string{}
	}
	if tags, ok := parseTagsFromJSONArray(extracted); ok {
		return normalizeTags(tags)
	}
	return []string{}
}

func parseTagsFromJSONArray(text string) ([]string, bool) {
	var tags []string
	if err := json.Unmarshal([]byte(text), &tags); err != nil {
		return nil, false
	}
	return tags, true
}

func extractFirstJSONArray(text string) (string, bool) {
	re := regexp.MustCompile(`\[[\s\S]*?\]`)
	m := re.FindString(text)
	if strings.TrimSpace(m) == "" {
		return "", false
	}
	return m, true
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		v := strings.TrimSpace(t)
		if v == "" {
			continue
		}
		if len([]rune(v)) > 12 {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
		if len(out) >= 10 {
			break
		}
	}
	return out
}
