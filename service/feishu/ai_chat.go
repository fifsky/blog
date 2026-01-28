package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"app/config"
	"app/pkg/mcp"

	"github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcardkit "github.com/larksuite/oapi-sdk-go/v3/service/cardkit/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// AIChat handles AI chat interactions with streaming card updates.
type AIChat struct {
	conf       *config.Config
	larkClient *lark.Client
	aiClient   openai.Client
	mcpManager *mcp.Manager
}

// NewAIChat creates a new AIChat instance.
func NewAIChat(conf *config.Config, larkClient *lark.Client) *AIChat {
	aiClient := openai.NewClient(
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

	return &AIChat{
		conf:       conf,
		larkClient: larkClient,
		aiClient:   aiClient,
		mcpManager: mcpManager,
	}
}

// CardUpdater manages card element updates with auto-incrementing sequence numbers.
type CardUpdater struct {
	cardID     string
	larkClient *lark.Client
	sequence   int64
}

// NewCardUpdater creates a new CardUpdater for the given card ID.
func NewCardUpdater(cardID string, larkClient *lark.Client) *CardUpdater {
	return &CardUpdater{
		cardID:     cardID,
		larkClient: larkClient,
		sequence:   0,
	}
}

// UpdateElement updates a card element's content with auto-incrementing sequence.
func (u *CardUpdater) UpdateElement(ctx context.Context, elementID, content string) error {
	seq := int(atomic.AddInt64(&u.sequence, 1))

	req := larkcardkit.NewContentCardElementReqBuilder().
		CardId(u.cardID).
		ElementId(elementID).
		Body(larkcardkit.NewContentCardElementReqBodyBuilder().
			Content(content).
			Uuid(uuid.NewString()).
			Sequence(seq).
			Build()).
		Build()

	resp, err := u.larkClient.Cardkit.V1.CardElement.Content(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("update card element failed: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// UpdateContent updates the content element.
func (u *CardUpdater) UpdateContent(ctx context.Context, content string) error {
	return u.UpdateElement(ctx, "content", content)
}

// UpdateTip updates the tip div element using Patch API (div elements don't support streaming updates).
func (u *CardUpdater) UpdateTip(ctx context.Context, tipText string) error {
	seq := int(atomic.AddInt64(&u.sequence, 1))

	// Use PartialElement to update the text.content property of the div
	partial := map[string]any{
		"text": map[string]any{
			"content":    tipText,
			"text_size":  "notation",
			"text_align": "left",
			"text_color": "grey",
		},
	}
	partialBytes, _ := json.Marshal(partial)
	partialElement := string(partialBytes)

	req := larkcardkit.NewPatchCardElementReqBuilder().
		CardId(u.cardID).
		ElementId("tip").
		Body(larkcardkit.NewPatchCardElementReqBodyBuilder().
			PartialElement(partialElement).
			Uuid(uuid.NewString()).
			Sequence(seq).
			Build()).
		Build()

	resp, err := u.larkClient.Cardkit.V1.CardElement.Patch(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("update tip element failed: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// HandleMessage processes a user message and responds with AI-generated content
// using streaming card updates for typewriter effect.
func (a *AIChat) HandleMessage(ctx context.Context, messageID, userMessage string) error {
	// Step 1: Create a streaming card with initial loading state
	cardID, msgID, err := a.createStreamingCard(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to create streaming card: %w", err)
	}

	// Create card updater with sequence management
	updater := NewCardUpdater(cardID, a.larkClient)

	// Step 2: Call AI streaming API and update card content progressively
	if err := a.streamAIResponse(ctx, updater, userMessage); err != nil {
		// Update card with error message and final tip
		_ = updater.UpdateContent(ctx, fmt.Sprintf("❌ AI 响应失败: %v", err))
		_ = updater.UpdateTip(ctx, "以上内容由 AI 生成，仅供参考")
		return fmt.Errorf("failed to stream AI response: %w", err)
	}

	fmt.Printf("[Feishu Bot] AI response completed for message %s, card %s\n", msgID, cardID)
	return nil
}

// buildCardJSON builds the initial card JSON with streaming mode enabled.
func buildCardJSON(content, tip string) string {
	card := map[string]any{
		"schema": "2.0",
		"config": map[string]any{
			"update_multi":   true,
			"streaming_mode": true,
			"streaming_config": map[string]any{
				"print_step":         map[string]any{"default": 1},
				"print_frequency_ms": map[string]any{"default": 70},
				"print_strategy":     "fast",
			},
		},
		"body": map[string]any{
			"direction": "vertical",
			"elements": []map[string]any{
				{
					"tag":        "markdown",
					"content":    content,
					"text_align": "left",
					"text_size":  "normal",
					"margin":     "0px 0px 0px 0px",
					"element_id": "content",
				},
				{
					"tag": "div",
					"text": map[string]any{
						"tag":        "plain_text",
						"content":    tip,
						"text_size":  "notation",
						"text_align": "left",
						"text_color": "grey",
					},
					"icon": map[string]any{
						"tag":   "standard_icon",
						"token": "robot_outlined",
						"color": "grey",
					},
					"margin":     "0px 0px 0px 0px",
					"element_id": "tip",
				},
			},
		},
	}
	jsonBytes, _ := json.Marshal(card)
	return string(jsonBytes)
}

// createStreamingCard creates a streaming card and returns the card ID.
func (a *AIChat) createStreamingCard(ctx context.Context, messageID string) (string, string, error) {
	// Build card with initial content and tip
	cardContent := buildCardJSON("", "努力回答中…")

	replyReq := larkim.NewReplyMessageReqBuilder().
		MessageId(messageID).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType("interactive").
			Content(cardContent).
			Build()).
		Build()

	replyResp, err := a.larkClient.Im.V1.Message.Reply(ctx, replyReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to send card message: %w", err)
	}

	if !replyResp.Success() {
		return "", "", fmt.Errorf("failed to send card message: code=%d, msg=%s", replyResp.Code, replyResp.Msg)
	}

	msgID := *replyResp.Data.MessageId
	fmt.Printf("[Feishu Bot] Sent card as message %s\n", msgID)

	// Get card ID from message ID for streaming updates
	cardID, err := a.getCardID(ctx, msgID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get card ID: %w", err)
	}

	fmt.Printf("[Feishu Bot] Got card ID %s for message %s\n", cardID, msgID)
	return cardID, msgID, nil
}

// getCardID retrieves the card ID from a message ID.
func (a *AIChat) getCardID(ctx context.Context, messageID string) (string, error) {
	req := larkcardkit.NewIdConvertCardReqBuilder().
		Body(larkcardkit.NewIdConvertCardReqBodyBuilder().
			MessageId(messageID).
			Build()).
		Build()

	resp, err := a.larkClient.Cardkit.V1.Card.IdConvert(ctx, req)
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", fmt.Errorf("id convert failed: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data == nil || resp.Data.CardId == nil {
		return "", fmt.Errorf("card ID not found in response")
	}

	return *resp.Data.CardId, nil
}

// getMCPTools returns the available tools from all MCP clients as OpenAI tool params
func (a *AIChat) getMCPTools(ctx context.Context) []openai.ChatCompletionToolUnionParam {
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
func (a *AIChat) executeTool(ctx context.Context, name string, arguments string) string {
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

// streamAIResponse calls the AI API and streams the response to the card.
func (a *AIChat) streamAIResponse(ctx context.Context, updater *CardUpdater, userMessage string) error {
	prompt := fmt.Sprintf(`You are a helpful assistant. Respond in the same language as the user's message.
Current Time: %s

When you encounter questions that you cannot answer directly, such as:
- Current events, news, or real-time information
- Recent research or developments in professional fields
- Specific facts you are uncertain about
- Information that may have changed after your training data cutoff

You should use the available tools to find accurate and up-to-date information.

请用简洁友好的方式回答用户的问题。`, time.Now().Format(time.DateTime))

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
		openai.UserMessage(userMessage),
	}

	// Get tools from all MCP clients
	tools := a.getMCPTools(ctx)

	aiReq := openai.ChatCompletionNewParams{
		Model:    a.conf.Common.AIModel,
		Messages: messages,
		Tools:    tools,
	}

	if strings.HasPrefix(a.conf.Common.AIModel, "doubao") {
		aiReq.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	var content strings.Builder
	updateInterval := 300 * time.Millisecond
	lastUpdate := time.Now()

	// Tool calling loop - handle tool calls until we get a final response
	for {
		stream := a.aiClient.Chat.Completions.NewStreaming(ctx, aiReq)
		acc := openai.ChatCompletionAccumulator{}
		hasToolCalls := false

		// Stream the response
		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			// Check if stream finished with tool_calls
			if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason == "tool_calls" {
				hasToolCalls = true
				break
			}

			// Stream content to card
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				content.WriteString(chunk.Choices[0].Delta.Content)

				// Update card at intervals to avoid rate limiting
				if time.Since(lastUpdate) >= updateInterval {
					if err := updater.UpdateContent(ctx, content.String()); err != nil {
						fmt.Printf("[Feishu Bot] Failed to update card content: %v\n", err)
					}
					lastUpdate = time.Now()
				}
			}
		}

		if err := stream.Err(); err != nil {
			return fmt.Errorf("AI stream error: %w", err)
		}

		// Handle tool calls if any
		if hasToolCalls && len(acc.Choices) > 0 && len(acc.Choices[0].Message.ToolCalls) > 0 {
			// Add assistant message with tool calls
			aiReq.Messages = append(aiReq.Messages, acc.Choices[0].Message.ToParam())

			// Execute all tool calls and add results
			for _, toolCall := range acc.Choices[0].Message.ToolCalls {
				// Get MCP display name for the tool
				mcpName := a.mcpManager.GetMCPDisplayName(toolCall.Function.Name)

				// Update tip to show tool calling status
				tipText := fmt.Sprintf("正在调用工具，%s", mcpName)
				if err := updater.UpdateTip(ctx, tipText); err != nil {
					fmt.Printf("[Feishu Bot] Failed to update tip: %v\n", err)
				}

				fmt.Printf("[Feishu Bot] Calling tool: %s (%s)\n", toolCall.Function.Name, mcpName)

				// Execute tool
				toolResult := a.executeTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)

				// Restore tip to default
				if err := updater.UpdateTip(ctx, "努力回答中…"); err != nil {
					fmt.Printf("[Feishu Bot] Failed to restore tip: %v\n", err)
				}

				aiReq.Messages = append(aiReq.Messages, openai.ToolMessage(toolResult, toolCall.ID))
			}

			// Continue the loop to get the next response
			continue
		}

		// No tool calls, we're done
		break
	}

	// Final update with complete content
	finalContent := content.String()
	if finalContent == "" {
		finalContent = "抱歉，我暂时无法回答您的问题。"
	}

	// Update content and tip
	if err := updater.UpdateContent(ctx, finalContent); err != nil {
		return fmt.Errorf("failed to update final content: %w", err)
	}

	if err := updater.UpdateTip(ctx, "以上内容由 AI 生成，仅供参考"); err != nil {
		return fmt.Errorf("failed to update tip: %w", err)
	}

	return nil
}
