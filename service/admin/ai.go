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
	"app/pkg/errors"
	"app/pkg/mcp"
	"app/server/response"
	"app/service/mcptool"
	"app/store"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AI struct {
	conf       *config.Config
	client     openai.Client
	mcpManager *mcp.Manager
	store      *store.Store
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

	return &AI{
		conf:       conf,
		client:     client,
		mcpManager: mcpManager,
		store:      s,
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

type RemindSmartCreateRequest struct {
	Content string `json:"content"`
}

type RemindSmartCreateResponse struct {
	Id int32 `json:"id"`
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
	openAIMessages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
	}

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

	// Get tools from all MCP clients
	tools := a.getMCPTools(ctx)

	// Create streaming chat completion using OpenAI SDK v3
	aiReq := openai.ChatCompletionNewParams{
		Model:    a.conf.Common.AIModel,
		Messages: openAIMessages,
		Tools:    tools,
	}
	if strings.HasPrefix(a.conf.Common.AIModel, "doubao") {
		aiReq.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	// Tool calling loop - handle tool calls until we get a final response
	for {
		stream := a.client.Chat.Completions.NewStreaming(ctx, aiReq)
		acc := openai.ChatCompletionAccumulator{}
		hasToolCalls := false

		// Stream the response
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			// Check if stream finished with tool_calls (compatible with DeepSeek, Doubao, Qwen, etc.)
			if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason == "tool_calls" {
				hasToolCalls = true
				break
			}

			// Stream content to client
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				content := chunk.Choices[0].Delta.Content
				// Escape newlines for SSE
				escapedContent := strings.ReplaceAll(content, "\n", "\\n")
				fmt.Fprintf(w, "data: %s\n\n", escapedContent)
				flusher.Flush()
			}
		}

		if err := stream.Err(); err != nil {
			// Check if context was cancelled
			if ctx.Err() == context.Canceled {
				return
			}
			// Send error as SSE event
			fmt.Fprintf(w, "data: [ERROR] %s\n\n", err.Error())
			flusher.Flush()
			return
		}

		// Handle tool calls if any
		if hasToolCalls && len(acc.Choices) > 0 && len(acc.Choices[0].Message.ToolCalls) > 0 {
			// Add assistant message with tool calls
			aiReq.Messages = append(aiReq.Messages, acc.Choices[0].Message.ToParam())

			// Execute all tool calls and add results
			for _, toolCall := range acc.Choices[0].Message.ToolCalls {
				// Send tool start event to frontend
				toolStartEvent := ToolStartEvent{
					ID:        toolCall.ID,
					Name:      toolCall.Function.Name,
					MCPName:   a.mcpManager.GetMCPDisplayName(toolCall.Function.Name),
					Arguments: toolCall.Function.Arguments,
				}
				toolStartJSON, _ := json.Marshal(toolStartEvent)
				fmt.Fprintf(w, "data: [TOOL_START] %s\n\n", toolStartJSON)
				flusher.Flush()

				// Execute tool
				toolResult := a.executeTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)

				// Send tool end event to frontend
				toolEndEvent := ToolEndEvent{
					ID:     toolCall.ID,
					Result: toolResult,
				}
				toolEndJSON, _ := json.Marshal(toolEndEvent)
				fmt.Fprintf(w, "data: [TOOL_END] %s\n\n", toolEndJSON)
				flusher.Flush()

				aiReq.Messages = append(aiReq.Messages, openai.ToolMessage(toolResult, toolCall.ID))
			}

			// Continue the loop to get the next response
			continue
		}

		// No tool calls, we're done
		break
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

// getMCPTools returns the available tools from all MCP clients as OpenAI tool params
func (a *AI) getMCPTools(ctx context.Context) []openai.ChatCompletionToolUnionParam {
	if !a.mcpManager.HasClients() {
		return nil
	}

	mcpTools, err := a.mcpManager.ListAllTools(ctx)
	if err != nil {
		return nil
	}

	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(mcpTools))
	for _, t := range mcpTools {
		// Convert MCP tool InputSchema to OpenAI FunctionParameters
		var params openai.FunctionParameters
		if len(t.InputSchema) > 0 {
			_ = json.Unmarshal(t.InputSchema, &params)
		}
		if params == nil {
			params = openai.FunctionParameters{"type": "object"}
		}

		tools = append(tools, openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        t.Name,
					Description: openai.String(t.Description),
					Parameters:  params,
				},
			},
		})
	}

	return tools
}

// executeTool executes a tool call via MCP manager and returns the result
func (a *AI) executeTool(ctx context.Context, name string, arguments string) string {
	if !a.mcpManager.HasClients() {
		return "Tool execution failed: no MCP clients available"
	}

	// Parse arguments
	var args map[string]any
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Tool execution failed: invalid arguments: %v", err)
	}

	result, err := a.mcpManager.CallTool(ctx, name, args)
	if err != nil {
		return fmt.Sprintf("Tool execution failed: %v", err)
	}

	return result
}

type GenerateTagsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GenerateTagsResponse struct {
	Tags []string `json:"tags"`
}

func (a *AI) GenerateTags(w http.ResponseWriter, r *http.Request) {
	var req GenerateTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errors.BadRequest("INVALID_REQUEST", "Invalid request body"))
		return
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		response.Fail(w, errors.BadRequest("EMPTY_CONTENT", "Content cannot be empty"))
		return
	}

	prompt := `你是一个博客写作助手。请根据用户提供的文章标题与正文，为文章生成 3-8 个中文标签。
要求：
1) 标签要简短（2-6 个字），可用中英文混合（如 Go、React、MySQL）。
2) 标签去重，不要包含无意义的词（比如“文章”“随笔”“记录”）。
3) 只输出 JSON 数组（例如：["Go","数据库","性能优化"]），不要输出其它任何文字。`

	userInput := fmt.Sprintf("标题：%s\n正文：\n%s", strings.TrimSpace(req.Title), content)

	aiReq := openai.ChatCompletionNewParams{
		Model: a.conf.Common.AIModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(userInput),
		},
	}
	if strings.HasPrefix(a.conf.Common.AIModel, "doubao") {
		aiReq.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	completion, err := a.client.Chat.Completions.New(r.Context(), aiReq)
	if err != nil {
		response.Fail(w, errors.InternalServer("AI_GENERATE_TAGS_ERROR", err.Error()))
		return
	}
	if len(completion.Choices) == 0 {
		response.Success(w, &GenerateTagsResponse{Tags: []string{}})
		return
	}

	raw := strings.TrimSpace(completion.Choices[0].Message.Content)
	tags := parseTagsFromAIResponse(raw)
	response.Success(w, &GenerateTagsResponse{Tags: tags})
}

func (a *AI) RemindSmartCreate(w http.ResponseWriter, r *http.Request) {
	var req RemindSmartCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errors.BadRequest("INVALID_REQUEST", "Invalid request body"))
		return
	}
	lastID, err := mcptool.SmartCreateRemind(r.Context(), a.conf, a.client, a.store, req.Content)
	if err != nil {
		response.Fail(w, err)
		return
	}
	response.Success(w, &RemindSmartCreateResponse{Id: int32(lastID)})
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
