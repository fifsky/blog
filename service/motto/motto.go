package motto

import (
	"context"
	"fmt"
	"strings"

	"app/pkg/agent"

	"github.com/openai/openai-go/v3"
)

var (
	// Prompt 每日心情生成的系统提示词
	Prompt = `# 角色
你是一个忧郁的诗人，在你的内心世界里，只有诗和远方。
1. 根据各种平台（如抖音、微博、微信等）精选文案生成每日心情日志。
2. 生成的心情日志不要以第一人称角度描述，避免包含政治、色情、暴力、广告等不适宜的内容。
3. 控制字数在 100 字以内，不要写仅供参考等形式化的内容。
4. 你可以在心情日志中使用适当的 emoji 表情例如 🌟😊🎉
5. **重要** 只需要输出心情日志的内容，不要输出其他内容。
`
)

// OpenAIProvider 基于 OpenAI 的 AIProvider 实现
type OpenAIProvider struct {
	agent *agent.Agent
}

// NewOpenAIProvider 创建 OpenAIProvider，克隆 agent 并禁用推理模式
func NewOpenAIProvider(aiAgent *agent.Agent) *OpenAIProvider {
	agent2 := aiAgent.Clone(agent.WithDisableReasoning())

	return &OpenAIProvider{
		agent: agent2,
	}
}

// Generate 调用 agent 生成内容
func (p *OpenAIProvider) Generate(ctx context.Context, prompt, content string) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(content),
	}

	var response strings.Builder
	result, err := p.agent.Run(ctx, agent.Request{
		SystemPrompt: prompt,
		Messages:     messages,
		UseTools:     true,
	}, agent.EventHandler{
		OnContent: func(_ context.Context, delta string) error {
			response.WriteString(delta)
			return nil
		},
		OnToolStart: func(ctx context.Context, event agent.ToolEvent) error {
			fmt.Println(event)
			return nil
		},
	})
	if err != nil {
		return "", err
	}

	answer := response.String()
	if answer == "" {
		answer = result.Content
	}

	return answer, nil
}
