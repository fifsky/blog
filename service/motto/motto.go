package motto

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"app/pkg/aiagent"
	"app/pkg/bark"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/robfig/cron/v3"
)

var (
	Prompt = `# 角色
根据我提供的日期，查询和我兴趣关注点相关的内容，生成一段符合意境的心情日志
1. **信息准确性守护者**：确保提供的信息准确无误。
2. 生成的心情日志必须符合我兴趣关注点相关的内容，不要以第一人称角度描述，避免包含政治、色情、暴力、广告等不适宜的内容。
3. 控制字数在 100 字以内，不要写仅供参考等形式化的内容。
4. **回答更生动活泼**：你可以在心情日志中使用适当的 emoji 表情来描述天气和心情，例如 🌟😊🎉
5. **重要** 只需要输出心情日志的内容，不要输出其他内容。

## 我兴趣关注点相关的内容
- 科技（人工智能、IT、编程）
- 文学
- 旅行
- 电影
- 音乐（民谣、流行、治愈）
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
	return &OpenAIProvider{
		agent: agent,
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
	store      *store.Store
	barkClient *bark.Client
	ai         AIProvider
}

func New(s *store.Store, barkClient *bark.Client, ai AIProvider) *Motto {
	return &Motto{
		store:      s,
		barkClient: barkClient,
		ai:         ai,
	}
}

func (m *Motto) Start(spec string) {
	c := cron.New()
	// 每天上午9点准时调用
	_, err := c.AddFunc(spec, func() {
		if err := m.GenerateDailyMotto(); err != nil {
			logger.Default().Error("generate daily motto error", slog.String("err", err.Error()))
		} else {
			logger.Default().Info("generate daily motto success")
		}
	})
	if err != nil {
		logger.Default().Error("motto cron add func error", slog.String("err", err.Error()))
		return
	}
	c.Start()
}

func (m *Motto) GenerateDailyMotto() error {
	logger.Default().Info("start generate daily motto")
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

	// 发送提醒
	if err := m.sendBark(content); err != nil {
		logger.Default().Error("motto request bark error", slog.String("err", err.Error()))
	}

	return nil
}

func (m *Motto) sendBark(content string) error {
	msg := bark.Message{
		Title: "每日一言",
		Body:  content,
		Badge: 1,
		Group: "Motto",
		Level: "timeSensitive",
	}

	return m.barkClient.Send(msg)
}
