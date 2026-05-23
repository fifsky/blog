// Package aiagent 提供 OpenAI 与 MCP 工具调用的通用编排流程。
package aiagent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"app/config"
	mcpclient "app/pkg/mcp"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type toolProvider interface {
	HasClients() bool
	ListAllTools(ctx context.Context) ([]mcpclient.Tool, error)
	CallTool(ctx context.Context, toolName string, arguments map[string]any) (string, error)
	GetMCPDisplayName(toolName string) string
}

type chatStream interface {
	Next() bool
	Current() openai.ChatCompletionChunk
	Err() error
}

type chatStreamFactory interface {
	NewStreaming(ctx context.Context, req openai.ChatCompletionNewParams) (chatStream, error)
}

type streamResult struct {
	acc              openai.ChatCompletionAccumulator
	reasoningContent string
}

type openAIStreamFactory struct {
	client openai.Client
}

func (f openAIStreamFactory) NewStreaming(ctx context.Context, req openai.ChatCompletionNewParams) (chatStream, error) {
	return f.client.Chat.Completions.NewStreaming(ctx, req), nil
}

// Agent 负责 OpenAI 初始化、MCP 工具绑定和工具调用循环。
type Agent struct {
	client           openai.Client
	model            string
	tools            toolProvider
	streamFactory    chatStreamFactory
	disableReasoning bool
	reasoningEffort  string
}

// Request 描述一次 AI 对话请求。
type Request struct {
	SystemPrompt string
	Messages     []openai.ChatCompletionMessageParamUnion
	UseTools     bool
}

// Result 描述 AI 编排完成后的最终结果。
type Result struct {
	Content  string
	Messages []openai.ChatCompletionMessageParamUnion
}

// ToolEvent 描述一次 MCP 工具调用事件。
type ToolEvent struct {
	ID        string
	Name      string
	MCPName   string
	Arguments string
	Result    string
}

// EventHandler 允许业务层处理流式文本与工具调用事件。
type EventHandler struct {
	OnContent   func(ctx context.Context, delta string) error
	OnReasoning func(ctx context.Context, delta string) error
	OnToolStart func(ctx context.Context, event ToolEvent) error
	OnToolEnd   func(ctx context.Context, event ToolEvent) error
}

// Option 定义 Agent 的配置选项
type Option func(*Agent)

// WithClient 设置 OpenAI 客户端
func WithClient(client openai.Client) Option {
	return func(a *Agent) {
		a.client = client
	}
}

// WithModel 设置使用的模型名称
func WithModel(model string) Option {
	return func(a *Agent) {
		a.model = model
	}
}

// WithDisableReasoning 设置是否禁用深度思考
func WithDisableReasoning() Option {
	return func(a *Agent) {
		a.disableReasoning = true
	}
}

// WithReasoningEffort 设置深度思考级别，默认为 high
func WithReasoningEffort(effort string) Option {
	return func(a *Agent) {
		if effort == "" {
			effort = "high"
		}
		a.reasoningEffort = effort
	}
}

// WithMCP 设置 MCP 配置
func WithMCP(mcp map[string]config.MCPConf) Option {
	return func(a *Agent) {
		manager := mcpclient.NewManager()
		for key, mcpConf := range mcp {
			if mcpConf.URL == "" {
				continue
			}
			displayName := mcpConf.Name
			if displayName == "" {
				displayName = key
			}
			manager.AddClient(key, displayName, mcpConf.URL, mcpConf.Token)
		}
		a.tools = manager
	}
}

// New 创建带配置的 Agent。
func New(opts ...Option) *Agent {
	a := &Agent{
		reasoningEffort: "high",
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// GetClient 返回当前使用的 OpenAI 客户端
func (a *Agent) GetClient() openai.Client {
	return a.client
}

// GetModel 返回当前使用的模型名称
func (a *Agent) GetModel() string {
	return a.model
}

// Run 执行流式对话，并在模型请求工具时调用 MCP 后继续生成。
func (a *Agent) Run(ctx context.Context, request Request, handler EventHandler) (Result, error) {
	streamFactory := a.streamFactory
	if streamFactory == nil {
		streamFactory = openAIStreamFactory{client: a.client}
	}

	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Messages)+1)
	if strings.TrimSpace(request.SystemPrompt) != "" {
		messages = append(messages, openai.SystemMessage(request.SystemPrompt))
	}
	messages = append(messages, request.Messages...)

	aiReq := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(a.model),
		Messages: messages,
	}
	if !a.disableReasoning {
		aiReq.ReasoningEffort = shared.ReasoningEffort(a.reasoningEffort)
		aiReq.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "enabled",
			},
		})
	}
	if request.UseTools {
		aiReq.Tools = a.buildTools(ctx)
	}

	var content strings.Builder
	generatedMessages := make([]openai.ChatCompletionMessageParamUnion, 0, 2)
	for {
		streamResult, err := a.runStream(ctx, streamFactory, aiReq, handler, &content)
		if err != nil {
			return Result{}, err
		}
		if len(streamResult.acc.Choices) == 0 {
			break
		}

		assistantMessage := assistantMessageWithReasoning(streamResult.acc.Choices[0].Message, streamResult.reasoningContent)
		generatedMessages = append(generatedMessages, assistantMessage)
		if len(streamResult.acc.Choices[0].Message.ToolCalls) == 0 {
			break
		}

		aiReq.Messages = append(aiReq.Messages, assistantMessage)
		for _, toolCall := range streamResult.acc.Choices[0].Message.ToolCalls {
			event := ToolEvent{
				ID:        toolCall.ID,
				Name:      toolCall.Function.Name,
				MCPName:   a.mcpDisplayName(toolCall.Function.Name),
				Arguments: toolCall.Function.Arguments,
			}
			if handler.OnToolStart != nil {
				if err := handler.OnToolStart(ctx, event); err != nil {
					return Result{}, err
				}
			}

			result := a.executeToolCall(ctx, toolCall)
			event.Result = result
			if handler.OnToolEnd != nil {
				if err := handler.OnToolEnd(ctx, event); err != nil {
					return Result{}, err
				}
			}
			toolMessage := openai.ToolMessage(result, toolCall.ID)
			generatedMessages = append(generatedMessages, toolMessage)
			aiReq.Messages = append(aiReq.Messages, toolMessage)
		}
	}

	return Result{Content: content.String(), Messages: generatedMessages}, nil
}

func (a *Agent) runStream(ctx context.Context, streamFactory chatStreamFactory, aiReq openai.ChatCompletionNewParams, handler EventHandler, content *strings.Builder) (streamResult, error) {
	stream, err := streamFactory.NewStreaming(ctx, aiReq)
	if err != nil {
		return streamResult{}, err
	}

	acc := openai.ChatCompletionAccumulator{}
	var reasoningContent strings.Builder
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) == 0 {
			continue
		}

		if reasoningDelta := reasoningContentFromDelta(chunk.Choices[0].Delta); reasoningDelta != "" {
			reasoningContent.WriteString(reasoningDelta)
			if handler.OnReasoning != nil {
				if err := handler.OnReasoning(ctx, reasoningDelta); err != nil {
					return streamResult{acc: acc, reasoningContent: reasoningContent.String()}, err
				}
			}
		}

		if chunk.Choices[0].Delta.Content != "" {
			delta := chunk.Choices[0].Delta.Content
			content.WriteString(delta)
			if handler.OnContent != nil {
				if err := handler.OnContent(ctx, delta); err != nil {
					return streamResult{acc: acc, reasoningContent: reasoningContent.String()}, err
				}
			}
		}

		if chunk.Choices[0].FinishReason == "tool_calls" {
			break
		}
	}
	if err := stream.Err(); err != nil {
		return streamResult{}, err
	}
	return streamResult{acc: acc, reasoningContent: reasoningContent.String()}, nil
}

func assistantMessageWithReasoning(message openai.ChatCompletionMessage, reasoningContent string) openai.ChatCompletionMessageParamUnion {
	assistantMessage := message.ToAssistantMessageParam()
	if reasoningContent != "" {
		// DeepSeek 工具调用后的后续请求必须完整回传 reasoning_content。
		setAssistantReasoningContent(&assistantMessage, reasoningContent)
	}
	return openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMessage}
}

// DecodeMessageParam 从 JSON 解码消息，并保留 OpenAI SDK 未建模的 DeepSeek 字段。
func DecodeMessageParam(rawMessage []byte) (openai.ChatCompletionMessageParamUnion, error) {
	var message openai.ChatCompletionMessageParamUnion
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		return message, err
	}

	var extra struct {
		Role             string `json:"role"`
		ReasoningContent string `json:"reasoning_content"`
	}
	if err := json.Unmarshal(rawMessage, &extra); err != nil {
		return message, nil
	}
	if extra.Role == "assistant" && extra.ReasoningContent != "" && message.OfAssistant != nil {
		setAssistantReasoningContent(message.OfAssistant, extra.ReasoningContent)
	}
	return message, nil
}

func setAssistantReasoningContent(message *openai.ChatCompletionAssistantMessageParam, reasoningContent string) {
	message.SetExtraFields(map[string]any{
		"reasoning_content": reasoningContent,
	})
}

func reasoningContentFromDelta(delta openai.ChatCompletionChunkChoiceDelta) string {
	field, ok := delta.JSON.ExtraFields["reasoning_content"]
	if ok && field.Valid() {
		var content string
		if err := json.Unmarshal([]byte(field.Raw()), &content); err == nil {
			return content
		}
	}

	var rawDelta struct {
		ReasoningContent string `json:"reasoning_content"`
	}
	if err := json.Unmarshal([]byte(delta.RawJSON()), &rawDelta); err != nil {
		return ""
	}
	return rawDelta.ReasoningContent
}

func (a *Agent) buildTools(ctx context.Context) []openai.ChatCompletionToolUnionParam {
	if a.tools == nil || !a.tools.HasClients() {
		return nil
	}

	mcpTools, err := a.tools.ListAllTools(ctx)
	if err != nil {
		return nil
	}

	tools := make([]openai.ChatCompletionToolUnionParam, 0, len(mcpTools))
	for _, t := range mcpTools {
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

func (a *Agent) executeToolCall(ctx context.Context, toolCall openai.ChatCompletionMessageToolCallUnion) string {
	if a.tools == nil || !a.tools.HasClients() {
		return "Tool execution failed: no MCP clients available"
	}

	var args map[string]any
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return fmt.Sprintf("Tool execution failed: invalid arguments: %v", err)
	}

	result, err := a.tools.CallTool(ctx, toolCall.Function.Name, args)
	if err != nil {
		return fmt.Sprintf("Tool execution failed: %v", err)
	}
	return result
}

func (a *Agent) mcpDisplayName(toolName string) string {
	if a.tools == nil {
		return ""
	}
	return a.tools.GetMCPDisplayName(toolName)
}
