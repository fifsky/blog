package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/config"
	"app/pkg/errors"
	"app/pkg/mcp"
	"app/server/response"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AI struct {
	conf      *config.Config
	client    openai.Client
	mcpClient *mcp.Client
}

func NewAI(conf *config.Config) *AI {
	client := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)

	var mcpClient *mcp.Client
	if conf.WebSearch.URL != "" {
		mcpClient = mcp.NewClient(conf.WebSearch.URL, conf.WebSearch.Token)
	}

	return &AI{
		conf:      conf,
		client:    client,
		mcpClient: mcpClient,
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

	// Get tools from MCP if available
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
				toolResult := a.executeTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
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

You should use the available search tool to find accurate and up-to-date information.`, time.Now().Format(time.DateTime))

	return basePrompt
}

// getMCPTools returns the available tools from MCP server as OpenAI tool params
func (a *AI) getMCPTools(ctx context.Context) []openai.ChatCompletionToolUnionParam {
	if a.mcpClient == nil {
		return nil
	}

	mcpTools, err := a.mcpClient.ListTools(ctx)
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

// executeTool executes a tool call via MCP and returns the result
func (a *AI) executeTool(ctx context.Context, name string, arguments string) string {
	if a.mcpClient == nil {
		return "Tool execution failed: MCP client not available"
	}

	// Parse arguments
	var args map[string]any
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("Tool execution failed: invalid arguments: %v", err)
	}

	result, err := a.mcpClient.CallTool(ctx, name, args)
	if err != nil {
		return fmt.Sprintf("Tool execution failed: %v", err)
	}

	return result
}
