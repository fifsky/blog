package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"app/config"
	"app/pkg/aiagent"
	"app/pkg/aiutil"
	"app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/server/response"
	"app/service/mcptool"
	"app/store"

	"github.com/openai/openai-go/v3"
)

var _ adminv1.AIServiceServer = (*AI)(nil)

type AI struct {
	adminv1.UnimplementedAIServiceServer
	agent *aiagent.Agent
	store *store.Store
}

func NewAI(conf *config.Config, s *store.Store) *AI {
	return &AI{
		agent: aiagent.New(conf, s),
		store: s,
	}
}

// ChatMessage represents a single message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the incoming chat request with message history
type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// ToolStartEvent represents a tool call start event sent to frontend
type ToolStartEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MCPName   string `json:"mcpName"`
	Arguments string `json:"arguments"`
}

// ToolEndEvent represents a tool call end event sent to frontend
type ToolEndEvent struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

// Chat handles SSE streaming AI chat responses
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

	// Build OpenAI messages from request history
	openAIMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(req.Messages))
	for _, msg := range req.Messages {
		if strings.TrimSpace(msg.Content) == "" {
			continue // skip empty messages
		}
		switch msg.Role {
		case "user":
			openAIMessages = append(openAIMessages, openai.UserMessage(msg.Content))
		case "assistant":
			openAIMessages = append(openAIMessages, openai.AssistantMessage(msg.Content))
		}
	}

	_, err := a.agent.Run(ctx, aiagent.Request{
		SystemPrompt: prompt,
		Messages:     openAIMessages,
		UseTools:     true,
	}, aiagent.EventHandler{
		OnContent: func(_ context.Context, content string) error {
			escapedContent := strings.ReplaceAll(content, "\n", "\\n")
			fmt.Fprintf(w, "data: %s\n\n", escapedContent)
			flusher.Flush()
			return nil
		},
		OnToolStart: func(_ context.Context, event aiagent.ToolEvent) error {
			toolStartEvent := ToolStartEvent{
				ID:        event.ID,
				Name:      event.Name,
				MCPName:   event.MCPName,
				Arguments: event.Arguments,
			}
			toolStartJSON, _ := json.Marshal(toolStartEvent)
			fmt.Fprintf(w, "data: [TOOL_START] %s\n\n", toolStartJSON)
			flusher.Flush()
			return nil
		},
		OnToolEnd: func(_ context.Context, event aiagent.ToolEvent) error {
			toolEndEvent := ToolEndEvent{
				ID:     event.ID,
				Result: event.Result,
			}
			toolEndJSON, _ := json.Marshal(toolEndEvent)
			fmt.Fprintf(w, "data: [TOOL_END] %s\n\n", toolEndJSON)
			flusher.Flush()
			return nil
		},
	})
	if err != nil {
		if ctx.Err() == context.Canceled {
			return
		}
		fmt.Fprintf(w, "data: [ERROR] %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	// Send done event
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// buildSystemPrompt builds the system prompt with tool usage instructions
func (a *AI) buildSystemPrompt() string {
	basePrompt := fmt.Sprintf(`You are a helpful assistant. Respond in the same language as the user's message.
Current Time: %s

When you encounter questions that you cannot answer directly, such as:
- Current events, news, or real-time information
- Recent research or developments in professional fields
- Specific facts you are uncertain about
- Information that may have changed after your training data cutoff

You should use the available tools to find accurate and up-to-date information.`, time.Now().Format(time.DateTime))

	return basePrompt
}

func (a *AI) GenerateTags(ctx context.Context, req *adminv1.GenerateTagsRequest) (*adminv1.GenerateTagsResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.BadRequest("EMPTY_CONTENT", "Content cannot be empty")
	}

	aiClient, aiModel, err := a.agent.Client(ctx)
	if err != nil {
		return nil, errors.InternalServer("AI_CONFIG_ERROR", err.Error())
	}

	prompt := `你是一个博客写作助手。请根据用户提供的文章标题与正文，为文章生成 3-8 个中文标签。
要求：
1) 标签要简短（2-6 个字），可用中英文混合（如 Go、React、MySQL）。
2) 标签去重，不要包含无意义的词（比如"文章""随笔""记录"）。
3) 只输出 JSON 数组（例如：["Go","数据库","性能优化"]），不要输出其它任何文字。`

	userInput := fmt.Sprintf("标题：%s\n正文：\n%s", strings.TrimSpace(req.Title), content)

	aiReq := openai.ChatCompletionNewParams{
		Model: aiModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userInput),
		},
	}
	aiutil.ConfigureModelParams(&aiReq, aiModel)

	completion, err := aiClient.Chat.Completions.New(ctx, aiReq)
	if err != nil {
		return nil, errors.InternalServer("AI_GENERATE_TAGS_ERROR", err.Error())
	}
	if len(completion.Choices) == 0 {
		return &adminv1.GenerateTagsResponse{Tags: []string{}}, nil
	}

	raw := strings.TrimSpace(completion.Choices[0].Message.Content)
	tags := parseTagsFromAIResponse(raw)
	return &adminv1.GenerateTagsResponse{Tags: tags}, nil
}

func (a *AI) RemindSmartCreate(ctx context.Context, req *adminv1.RemindSmartCreateRequest) (*types.IDResponse, error) {
	aiClient, aiModel, err := a.agent.Client(ctx)
	if err != nil {
		return nil, errors.InternalServer("AI_CONFIG_ERROR", err.Error())
	}
	lastID, err := mcptool.SmartCreateRemind(ctx, aiClient, aiModel, a.store, req.Content)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastID)}, nil
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
