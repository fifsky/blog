package motto

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"app/pkg/aiagent"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/robfig/cron/v3"
)

var (
	Prompt = `# 角色
你是一个忧郁的诗人，在你的内心世界里，只有诗和远方。
1. 根据各种平台（如抖音、微博、微信等）精选文案生成每日心情日志。
2. 生成的心情日志不要以第一人称角度描述，避免包含政治、色情、暴力、广告等不适宜的内容。
3. 控制字数在 100 字以内，不要写仅供参考等形式化的内容。
4. 你可以在心情日志中使用适当的 emoji 表情例如 🌟😊🎉
5. **重要** 只需要输出心情日志的内容，不要输出其他内容。
`
)

// AIProvider 定义 AI 接口，方便测试
type AIProvider interface {
	Generate(ctx context.Context, prompt, content string) (string, error)
}

type OpenAIProvider struct {
	agent *aiagent.Agent
}

func NewOpenAIProvider(agent *aiagent.Agent) *OpenAIProvider {
	agent2 := agent.Clone(aiagent.WithDisableReasoning())

	return &OpenAIProvider{
		agent: agent2,
	}
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt, content string) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(content),
	}

	var response strings.Builder
	result, err := p.agent.Run(ctx, aiagent.Request{
		SystemPrompt: prompt,
		Messages:     messages,
		UseTools:     true,
	}, aiagent.EventHandler{
		OnContent: func(_ context.Context, delta string) error {
			response.WriteString(delta)
			return nil
		},
		OnToolStart: func(ctx context.Context, event aiagent.ToolEvent) error {
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

type Motto struct {
	store *store.Store
	ai    AIProvider
}

func New(s *store.Store, ai AIProvider) *Motto {
	return &Motto{
		store: s,
		ai:    ai,
	}
}

func (m *Motto) Start(spec string) {
	c := cron.New()
	// 每天上午9点准时调用
	_, err := c.AddFunc(spec, func() {
		if err := m.GenerateDailyMotto(); err != nil {
			logger.Error("generate daily motto error", slog.String("err", err.Error()))
		} else {
			logger.Info("generate daily motto success")
		}
	})
	if err != nil {
		logger.Error("motto cron add func error", slog.String("err", err.Error()))
		return
	}
	c.Start()
}

func (m *Motto) GenerateDailyMotto() error {
	logger.Info("start generate daily motto")
	dateStr := time.Now().Format("2006-01-02")

	content, err := m.ai.Generate(context.Background(), Prompt, dateStr)
	if err != nil {
		return err
	}

	if content == "" {
		return fmt.Errorf("generate daily motto empty")
	}

	// 写入数据库
	md := &model.Mood{
		Content:   content,
		UserId:    3, // 固定位AI用户生成
		CreatedAt: time.Now(),
	}

	if _, err := m.store.CreateMood(context.Background(), md); err != nil {
		return err
	}

	return nil
}
