package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app/config"
	"app/pkg/aiutil"
	"app/pkg/errors"
	"app/pkg/mcp"
	"app/pkg/skill"
	"app/pkg/tool"
	"app/pkg/tool/bash"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/server/response"
	"app/service/mcptool"
	"app/store"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

var _ adminv1.AIServiceServer = (*AI)(nil)

type AI struct {
	adminv1.UnimplementedAIServiceServer
	conf         *config.Config
	client       openai.Client
	mcpManager   *mcp.Manager
	skillManager *skill.Manager
	resolver     tool.Resolver
	store        *store.Store
}

func NewAI(conf *config.Config, s *store.Store) *AI {
	client := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)

	// Create MCP manager with all configured MCP clients
	mcpManager := mcp.NewManager()
	for key, mcpConf := range conf.MCP {
		if mcpConf.URL != "" {
			displayName := mcpConf.Name
			if displayName == "" {
				displayName = key
			}
			mcpManager.AddClient(key, displayName, mcpConf.URL, mcpConf.Token)
		}
	}

	skillManager := skill.NewManager(conf.Common.SkillsPath)
	_ = skillManager.Load()

	return &AI{
		conf:         conf,
		client:       client,
		mcpManager:   mcpManager,
		skillManager: skillManager,
		resolver:     tool.Resolvers{skillManager, mcpManager, tool.WrapSingleResolver(bash.New())},
		store:        s,
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

// ThinkingEvent represents a reasoning/thinking state update
type ThinkingEvent struct {
	Content  string `json:"content"`
	Thinking bool   `json:"thinking"`
	Duration string `json:"duration,omitempty"`
}

func (a *AI) sendEvent(w http.ResponseWriter, format string, args ...any) {
	fmt.Fprintf(w, format+"\n\n", args...)
	w.(http.Flusher).Flush()
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

	if _, ok := w.(http.Flusher); !ok {
		response.Fail(w, errors.InternalServer("STREAMING_ERROR", "Streaming not supported"))
		return
	}

	ctx := r.Context()

	// Build OpenAI messages from request history
	openAIMessages := a.buildMessages(req.Messages)

	// Get tools from resolver
	toolsList, toolsParams := a.getTools(ctx)

	// Create streaming chat completion using OpenAI SDK v3
	aiReq := a.buildAIRequest(openAIMessages, toolsParams)

	// Tool calling loop - handle tool calls until we get a final response
	for {
		stream := a.client.Chat.Completions.NewStreaming(ctx, aiReq)

		hasToolCalls, acc, err := a.processStream(ctx, w, stream)
		if err != nil {
			if err == context.Canceled {
				return
			}
			// Send error as SSE event
			a.sendEvent(w, "data: [ERROR] %s", err.Error())
			return
		}

		// Handle tool calls if any
		if hasToolCalls && len(acc.Choices) > 0 && len(acc.Choices[0].Message.ToolCalls) > 0 {
			a.handleToolCalls(ctx, w, &aiReq, acc.Choices[0].Message, toolsList)
			continue
		}

		// No tool calls, we're done
		break
	}

	// Send done event
	a.sendEvent(w, "data: [DONE]")
}

func (a *AI) buildMessages(reqMessages []ChatMessage) []openai.ChatCompletionMessageParamUnion {
	prompt := a.buildSystemPrompt()
	openAIMessages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
	}

	for _, msg := range reqMessages {
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
	return openAIMessages
}

func (a *AI) buildAIRequest(messages []openai.ChatCompletionMessageParamUnion, toolsParams []openai.ChatCompletionToolUnionParam) openai.ChatCompletionNewParams {
	aiReq := openai.ChatCompletionNewParams{
		Model:    a.conf.Common.AIModel,
		Messages: messages,
	}
	if len(toolsParams) > 0 {
		aiReq.Tools = toolsParams
	}
	aiutil.ConfigureModelParams(&aiReq, a.conf.Common.AIModel)
	return aiReq
}

func (a *AI) processStream(
	ctx context.Context,
	w http.ResponseWriter,
	stream *ssestream.Stream[openai.ChatCompletionChunk],
) (bool, openai.ChatCompletionAccumulator, error) {
	acc := openai.ChatCompletionAccumulator{}
	hasToolCalls := false

	var thinkStartTime time.Time
	isThinking := false

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason == "tool_calls" {
			hasToolCalls = true
			break
		}

		if len(chunk.Choices) > 0 {
			a.processChunkThinking(chunk, w, &isThinking, &thinkStartTime)

			// Stream regular content to client
			if chunk.Choices[0].Delta.Content != "" {
				content := chunk.Choices[0].Delta.Content
				escapedContent := strings.ReplaceAll(content, "\n", "\\n")
				a.sendEvent(w, "data: %s", escapedContent)
			}
		}
	}

	return hasToolCalls, acc, stream.Err()
}

func (a *AI) processChunkThinking(
	chunk openai.ChatCompletionChunk,
	w http.ResponseWriter,
	isThinking *bool,
	thinkStartTime *time.Time,
) {
	reasoningContent, has := chunk.Choices[0].Delta.JSON.ExtraFields["reasoning_content"]
	var rc string
	if has {
		if v, err := strconv.Unquote(reasoningContent.Raw()); err == nil {
			rc = v
		} else {
			rc = string(reasoningContent.Raw())
		}
	}

	if rc != "" {
		if !*isThinking {
			*isThinking = true
			*thinkStartTime = time.Now()
		}

		escapedRc := strings.ReplaceAll(rc, "\n", "\\n")
		thinkEv := ThinkingEvent{
			Content:  escapedRc,
			Thinking: true,
		}
		evJSON, _ := json.Marshal(thinkEv)
		a.sendEvent(w, "data: [THINKING] %s", evJSON)
	} else if *isThinking {
		*isThinking = false
		thinkEv := ThinkingEvent{
			Thinking: false,
			Duration: fmt.Sprintf("%.1f", time.Since(*thinkStartTime).Seconds()),
		}
		evJSON, _ := json.Marshal(thinkEv)
		a.sendEvent(w, "data: [THINKING] %s", evJSON)
	}
}

func (a *AI) handleToolCalls(
	ctx context.Context,
	w http.ResponseWriter,
	aiReq *openai.ChatCompletionNewParams,
	message openai.ChatCompletionMessage,
	toolsList []tool.Tool,
) {
	// Add assistant message with tool calls
	aiReq.Messages = append(aiReq.Messages, message.ToParam())

	for _, toolCall := range message.ToolCalls {
		mcpName := a.mcpManager.GetMCPDisplayName(toolCall.Function.Name)

		// Send tool start event to frontend
		toolStartEvent := ToolStartEvent{
			ID:        toolCall.ID,
			Name:      toolCall.Function.Name,
			MCPName:   mcpName,
			Arguments: toolCall.Function.Arguments,
		}
		toolStartJSON, _ := json.Marshal(toolStartEvent)
		a.sendEvent(w, "data: [TOOL_START] %s", toolStartJSON)

		// Execute tool using toolsList
		var toolResult string
		var found bool
		for _, t := range toolsList {
			if t.Name() == toolCall.Function.Name {
				res, err := t.Handle(ctx, toolCall.Function.Arguments)
				if err != nil {
					toolResult = fmt.Sprintf("Tool execution failed: %v", err)
				} else {
					toolResult = res
				}
				found = true
				break
			}
		}
		if !found {
			toolResult = fmt.Sprintf("Tool execution failed: tool %s not found", toolCall.Function.Name)
		}

		// Send tool end event to frontend
		toolEndEvent := ToolEndEvent{
			ID:     toolCall.ID,
			Result: toolResult,
		}
		toolEndJSON, _ := json.Marshal(toolEndEvent)
		a.sendEvent(w, "data: [TOOL_END] %s", toolEndJSON)

		aiReq.Messages = append(aiReq.Messages, openai.ToolMessage(toolResult, toolCall.ID))
	}
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

You should use the available tools to find accurate and up-to-date information.

%s`, time.Now().Format(time.DateTime), skill.SkillsPrompt(a.skillManager.GetSkills()))

	return basePrompt
}

// getTools returns the available tools as OpenAI tool params
func (a *AI) getTools(ctx context.Context) ([]tool.Tool, []openai.ChatCompletionToolUnionParam) {
	allTools, err := a.resolver.Resolve(ctx)
	if err != nil || len(allTools) == 0 {
		return nil, nil
	}

	params := make([]openai.ChatCompletionToolUnionParam, 0, len(allTools))
	for _, t := range allTools {
		var p openai.FunctionParameters
		if schema := t.InputSchema(); len(schema) > 0 {
			_ = json.Unmarshal(schema, &p)
		}
		if p == nil {
			p = openai.FunctionParameters{"type": "object"}
		}

		params = append(params, openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        t.Name(),
					Description: openai.String(t.Description()),
					Parameters:  p,
				},
			},
		})
	}

	return allTools, params
}

func (a *AI) GenerateTags(ctx context.Context, req *adminv1.GenerateTagsRequest) (*adminv1.GenerateTagsResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.BadRequest("EMPTY_CONTENT", "Content cannot be empty")
	}

	prompt := `你是一个博客写作助手。请根据用户提供的文章标题与正文，为文章生成 3-8 个中文标签。
要求：
1) 标签要简短（2-6 个字），可用中英文混合（如 Go、React、MySQL）。
2) 标签去重，不要包含无意义的词（比如"文章""随笔""记录"）。
3) 只输出 JSON 数组（例如：["Go","数据库","性能优化"]），不要输出其它任何文字。`

	userInput := fmt.Sprintf("标题：%s\n正文：\n%s", strings.TrimSpace(req.Title), content)

	aiReq := openai.ChatCompletionNewParams{
		Model: a.conf.Common.AIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userInput),
		},
	}
	aiutil.ConfigureModelParams(&aiReq, a.conf.Common.AIModel)

	completion, err := a.client.Chat.Completions.New(ctx, aiReq)
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
	lastID, err := mcptool.SmartCreateRemind(ctx, a.conf, a.client, a.store, req.Content)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastID)}, nil
}

func (a *AI) ListSkills(ctx context.Context, req *adminv1.ListSkillsRequest) (*adminv1.ListSkillsResponse, error) {
	skills := a.skillManager.GetSkills()
	res := &adminv1.ListSkillsResponse{
		Skills: make([]*adminv1.SkillInfo, 0, len(skills)),
	}
	for name, s := range skills {
		res.Skills = append(res.Skills, &adminv1.SkillInfo{
			Name:        name,
			Description: s.Description(),
		})
	}
	return res, nil
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
