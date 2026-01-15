package motto

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"app/config"
	"app/model"
	"app/pkg/bark"
	"app/store"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
	"github.com/robfig/cron/v3"
)

// AIProvider 定义 AI 接口，方便测试
type AIProvider interface {
	Generate(ctx context.Context, prompt, content string) (string, error)
}

type OpenAIProvider struct {
	client  *openai.Client
	model   string
	history []openai.ChatCompletionMessageParamUnion
	mu      sync.Mutex
}

func NewOpenAIProvider(token, endpoint, model string) *OpenAIProvider {
	client := openai.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL(endpoint),
	)
	return &OpenAIProvider{
		client: &client,
		model:  model,
	}
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt, content string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 构造消息上下文：系统提示词 + 历史记录 + 当前用户输入
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(prompt),
	}
	messages = append(messages, p.history...)
	messages = append(messages, openai.UserMessage(content))

	completion, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    shared.ChatModel(p.model),
	})

	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		response := completion.Choices[0].Message.Content
		// 记录历史消息：用户输入和AI输出
		p.history = append(p.history, openai.UserMessage(content))
		p.history = append(p.history, openai.AssistantMessage(response))
		return response, nil
	}
	return "", nil
}

type Motto struct {
	store      *store.Store
	conf       *config.Config
	barkClient *bark.Client
	ai         AIProvider
}

func New(s *store.Store, conf *config.Config, barkClient *bark.Client, ai AIProvider) *Motto {
	return &Motto{
		store:      s,
		conf:       conf,
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
	prompt := "你的任务是生成每日一言，用户告知你日期，你来生成一句名言，你只需要输出名言即可，不要输出其他内容"
	dateStr := time.Now().Format("2006-01-02")

	content, err := m.ai.Generate(context.Background(), prompt, dateStr)
	if err != nil {
		return err
	}

	if content == "" {
		return fmt.Errorf("generate daily motto empty")
	}

	// 写入数据库
	md := &model.Mood{
		Content:   content,
		UserId:    1, // 默认为管理员ID，假设为1
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
	}

	return m.barkClient.Send(msg)
}
