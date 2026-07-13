package openapi

import (
	"context"
	"fmt"
	"html"
	"time"

	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
)

var _ apiv1.GuestbookServiceHTTPServer = (*Guestbook)(nil)

// ContentModerator 内容审核接口
type ContentModerator interface {
	Moderate(ctx context.Context, content string) error
}

type Guestbook struct {
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

func NewGuestbook(s *store.Store, opts ...GuestbookOption) *Guestbook {
	g := &Guestbook{
		store: s,
	}

	// 应用选项
	for _, opt := range opts {
		opt(g)
	}

	// 如果没有设置自定义审核器，使用默认的 AI 审核器
	if g.moderator == nil {
		g.moderator = NewAIModerator(s)
	}

	return g
}

func (g *Guestbook) List(ctx context.Context, req *apiv1.GuestbookListRequest) (*apiv1.GuestbookListResponse, error) {
	num := 10
	guestbooks, err := g.store.ListGuestbook(ctx, req.GetKeyword(), int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}

	items := lo.Map(guestbooks, func(gb model.Guestbook, _ int) *apiv1.GuestbookItem {
		return apiv1.GuestbookItem_builder{Id: int32(gb.Id),
			Name:      gb.Name,
			Content:   gb.Content,
			Ip:        gb.Ip,
			CreatedAt: gb.CreatedAt.Format(time.DateTime),
			Top:       int32(gb.Top)}.Build()
	})

	total, err := g.store.CountGuestbookTotal(ctx, req.GetKeyword())
	if err != nil {
		return nil, err
	}

	return apiv1.GuestbookListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (g *Guestbook) Create(ctx context.Context, req *apiv1.GuestbookCreateRequest) (*apiv1.GuestbookCreateResponse, error) {
	// 内容审核
	if g.moderator != nil {
		if err := g.moderator.Moderate(ctx, fmt.Sprintf("%s %s", req.GetName(), req.GetContent())); err != nil {
			return nil, err
		}
	}

	// 获取客户端 IP
	ip := middleware.ClientIPFromContext(ctx)

	gb := &model.Guestbook{
		Name:      html.EscapeString(req.GetName()),
		Content:   html.EscapeString(req.GetContent()),
		Ip:        ip,
		CreatedAt: time.Now(),
	}

	id, err := g.store.CreateGuestbook(ctx, gb)
	if err != nil {
		return nil, err
	}

	return apiv1.GuestbookCreateResponse_builder{Id: int32(id)}.Build(),
		nil
}
