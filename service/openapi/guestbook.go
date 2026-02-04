package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"app/config"
	"app/pkg/aiutil"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/store"
	"app/store/model"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var _ apiv1.GuestbookServiceServer = (*Guestbook)(nil)

// ContentModerator 内容审核接口
type ContentModerator interface {
	Moderate(ctx context.Context, content string) error
}

type Guestbook struct {
	apiv1.UnimplementedGuestbookServiceServer
	store     *store.Store
	moderator ContentModerator
}

// GuestbookOption 用于配置 Guestbook 的选项函数
type GuestbookOption func(*Guestbook)

// WithModerator 设置自定义的内容审核器
func WithModerator(m ContentModerator) GuestbookOption {
	return func(g *Guestbook) {
		g.moderator = m
	}
}

func NewGuestbook(s *store.Store, conf *config.Config, opts ...GuestbookOption) *Guestbook {
	g := &Guestbook{
		store: s,
	}

	// 应用选项
	for _, opt := range opts {
		opt(g)
	}

	// 如果没有设置自定义审核器，使用默认的 AI 审核器
	if g.moderator == nil && conf != nil {
		g.moderator = NewAIModerator(conf)
	}

	return g
}

func (g *Guestbook) List(ctx context.Context, req *apiv1.GuestbookListRequest) (*apiv1.GuestbookListResponse, error) {
	num := 10
	guestbooks, err := g.store.ListGuestbook(ctx, req.Keyword, int(req.Page), num)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.GuestbookItem, 0, len(guestbooks))
	for _, gb := range guestbooks {
		item := &apiv1.GuestbookItem{
			Id:        int32(gb.Id),
			Name:      gb.Name,
			Content:   gb.Content,
			Ip:        gb.Ip,
			CreatedAt: gb.CreatedAt.Format(time.DateTime),
			Top:       int32(gb.Top),
		}
		items = append(items, item)
	}

	total, err := g.store.CountGuestbookTotal(ctx, req.Keyword)
	if err != nil {
		return nil, err
	}

	return &apiv1.GuestbookListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (g *Guestbook) Create(ctx context.Context, req *apiv1.GuestbookCreateRequest) (*apiv1.GuestbookCreateResponse, error) {
	// 内容审核
	if g.moderator != nil {
		if err := g.moderator.Moderate(ctx, fmt.Sprintf("%s %s", req.Name, req.Content)); err != nil {
			return nil, err
		}
	}

	// 获取客户端 IP
	ip := middleware.ClientIPFromContext(ctx)

	gb := &model.Guestbook{
		Name:      html.EscapeString(req.Name),
		Content:   html.EscapeString(req.Content),
		Ip:        ip,
		CreatedAt: time.Now(),
	}

	id, err := g.store.CreateGuestbook(ctx, gb)
	if err != nil {
		return nil, err
	}

	return &apiv1.GuestbookCreateResponse{
		Id: int32(id),
	}, nil
}

// AIModerator 基于 AI 的内容审核器
type AIModerator struct {
	aiClient openai.Client
	aiModel  string
}

// NewAIModerator 创建 AI 内容审核器
func NewAIModerator(conf *config.Config) *AIModerator {
	client := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)

	return &AIModerator{
		aiClient: client,
		aiModel:  conf.Common.AIModel,
	}
}

// Moderate 使用 AI 对内容进行审核
func (m *AIModerator) Moderate(ctx context.Context, content string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	prompt := `你是一个内容安全审核助手。请审核以下用户提交的留言内容，判断是否包含以下任何一种违规内容：
1. 色情、性感内容
2. 涉政敏感内容
3. 暴力恐怖内容
4. 违禁品相关内容
5. 宗教极端内容
6. 引流广告、垃圾广告
7. 辱骂、歧视、仇恨言论
8. 其他不良内容

请只回复 JSON 格式：
- 如果内容合规，回复：{"pass": true}
- 如果内容违规，回复：{"pass": false, "reason": "违规原因简述"}

不要输出任何其他内容。`

	aiReq := openai.ChatCompletionNewParams{
		Model: m.aiModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(content),
		},
	}
	aiutil.ConfigureModelParams(&aiReq, m.aiModel)

	completion, err := m.aiClient.Chat.Completions.New(ctx, aiReq)
	if err != nil {
		return errors.InternalServer("CONTENT_MODERATION_ERROR", "内容审核服务异常")
	}

	if len(completion.Choices) == 0 {
		return errors.InternalServer("CONTENT_MODERATION_ERROR", "内容审核服务无响应")
	}

	result := strings.TrimSpace(completion.Choices[0].Message.Content)

	// 解析 AI 返回的审核结果
	var moderationResult struct {
		Pass   bool   `json:"pass"`
		Reason string `json:"reason"`
	}

	if err := json.Unmarshal([]byte(result), &moderationResult); err != nil {
		return errors.InternalServer("CONTENT_MODERATION_ERROR", "内容审核结果解析失败")
	}

	if !moderationResult.Pass {
		reason := moderationResult.Reason
		if reason == "" {
			reason = "内容包含违规信息"
		}
		return errors.BadRequest("CONTENT_MODERATION_FAILED", reason)
	}

	return nil
}
