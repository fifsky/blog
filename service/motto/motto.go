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
	"app/pkg/doubao"
	"app/store"

	"github.com/goapt/logger"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
	"github.com/robfig/cron/v3"
)

var (
	prompt = `# 角色
每天自动根据用户提供的城市和日期查询天气（如：暴雨、雾霾、晚霞）生成一段符合意境的诗句或短语
1. **信息准确性守护者**：确保提供的信息准确无误。
2. 生成的诗句或短语必须符合意境，不要局限于城市信息。
3. **回答更生动活泼**：请在模型的回复中使用适当的 emoji 标签作为天气和心情的表示 🌟😊🎉"
`
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

type DoubaoProvider struct {
	client *doubao.Client
	model  string
}

func NewDoubaoProvider(apiKey, model string) *DoubaoProvider {
	return &DoubaoProvider{
		client: doubao.NewClient(apiKey),
		model:  model,
	}
}

func (p *DoubaoProvider) Generate(ctx context.Context, prompt, content string) (string, error) {
	resp, err := p.client.CreateChatCompletion(ctx, &doubao.ChatRequest{
		Model: p.model,
		Tools: []doubao.Tool{
			{
				Type:       "web_search",
				MaxKeyword: 2,
				Limit:      2,
			},
		},
		MaxToolCalls: 1,
		Thinking: &doubao.Thinking{
			Type: "disabled",
		},
		Input: []doubao.Message{
			{
				Role: "system",
				Content: []doubao.MessageContent{
					{
						Type: "input_text",
						Text: prompt,
					},
				},
			},
			{
				Role: "user",
				Content: []doubao.MessageContent{
					{
						Type: "input_text",
						Text: "城市：上海, 日期：2026-01-17",
					},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	for _, choice := range resp.Output {
		if choice.Type == "message" && len(choice.Content) > 0 {
			return choice.Content[0].Text, nil
		}
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
