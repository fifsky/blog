// Package aiagent 提供 OpenAI 与 MCP 工具调用的通用编排流程。
package aiagent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"app/config"
	"app/pkg/aiutil"
	mcpclient "app/pkg/mcp"
	"app/store"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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

type openAIStreamFactory struct {
	client openai.Client
}

func (f openAIStreamFactory) NewStreaming(ctx context.Context, req openai.ChatCompletionNewParams) (chatStream, error) {
	return f.client.Chat.Completions.NewStreaming(ctx, req), nil
}

// Agent 负责 OpenAI 初始化、MCP 工具绑定和工具调用循环。
type Agent struct {
	store          *store.Store
	tools          toolProvider
	clientProvider func(ctx context.Context) (openai.Client, string, error)
	streamFactory  chatStreamFactory
}

// Request 描述一次 AI 对话请求。
type Request struct {
	SystemPrompt string
	Messages     []openai.ChatCompletionMessageParamUnion
	UseTools     bool
}

// Result 描述 AI 编排完成后的最终结果。
type Result struct {
	Content string
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
	OnToolStart func(ctx context.Context, event ToolEvent) error
	OnToolEnd   func(ctx context.Context, event ToolEvent) error
}

// New 创建带数据库 AI 配置与 MCP 配置的 Agent。
func New(conf *config.Config, s *store.Store) *Agent {
	manager := mcpclient.NewManager()
	for key, mcpConf := range conf.MCP {
		if mcpConf.URL == "" {
			continue
		}
		displayName := mcpConf.Name
		if displayName == "" {
			displayName = key
		}
		manager.AddClient(key, displayName, mcpConf.URL, mcpConf.Token)
	}

	return &Agent{
		store: s,
		tools: manager,
	}
}

// Client 按需创建 OpenAI client，优先使用数据库配置。
func (a *Agent) Client(ctx context.Context) (openai.Client, string, error) {
	if a.clientProvider != nil {
		return a.clientProvider(ctx)
	}
	if a.store == nil {
		return openai.Client{}, "", fmt.Errorf("ai store is nil")
	}

	aiCfg := a.store.GetAIConfig(ctx)
	if aiCfg.Token == "" {
		return openai.Client{}, "", fmt.Errorf("ai token is empty")
	}

	client := openai.NewClient(
		option.WithAPIKey(aiCfg.Token),
		option.WithBaseURL(aiCfg.Endpoint),
	)
	return client, aiCfg.Model, nil
}

// Run 执行流式对话，并在模型请求工具时调用 MCP 后继续生成。
func (a *Agent) Run(ctx context.Context, request Request, handler EventHandler) (Result, error) {
	client, model, err := a.Client(ctx)
	if err != nil {
		return Result{}, err
	}

	streamFactory := a.streamFactory
	if streamFactory == nil {
		streamFactory = openAIStreamFactory{client: client}
	}

	messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Messages)+1)
	if strings.TrimSpace(request.SystemPrompt) != "" {
		messages = append(messages, openai.SystemMessage(request.SystemPrompt))
	}
	messages = append(messages, request.Messages...)

	aiReq := openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
	}
	if request.UseTools {
		aiReq.Tools = a.buildTools(ctx)
	}
	aiutil.ConfigureModelParams(&aiReq, model)

	var content strings.Builder
	for {
		acc, err := a.runStream(ctx, streamFactory, aiReq, handler, &content)
		if err != nil {
			return Result{}, err
		}
		if len(acc.Choices) == 0 || len(acc.Choices[0].Message.ToolCalls) == 0 {
			break
		}

		aiReq.Messages = append(aiReq.Messages, acc.Choices[0].Message.ToParam())
		for _, toolCall := range acc.Choices[0].Message.ToolCalls {
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
			aiReq.Messages = append(aiReq.Messages, openai.ToolMessage(result, toolCall.ID))
		}
	}

	return Result{Content: content.String()}, nil
}

func (a *Agent) runStream(ctx context.Context, streamFactory chatStreamFactory, aiReq openai.ChatCompletionNewParams, handler EventHandler, content *strings.Builder) (openai.ChatCompletionAccumulator, error) {
	stream, err := streamFactory.NewStreaming(ctx, aiReq)
	if err != nil {
		return openai.ChatCompletionAccumulator{}, err
	}

	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			delta := chunk.Choices[0].Delta.Content
			content.WriteString(delta)
			if handler.OnContent != nil {
				if err := handler.OnContent(ctx, delta); err != nil {
					return acc, err
				}
			}
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason == "tool_calls" {
			break
		}
	}
	if err := stream.Err(); err != nil {
		return acc, err
	}
	return acc, nil
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
