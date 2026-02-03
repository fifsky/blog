package openapi

import (
	"context"
	"html"
	"time"

	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/store"
	"app/store/model"
)

var _ apiv1.GuestbookServiceServer = (*Guestbook)(nil)

type Guestbook struct {
	apiv1.UnimplementedGuestbookServiceServer
	store *store.Store
}

func NewGuestbook(s *store.Store) *Guestbook {
	return &Guestbook{store: s}
}

func (g *Guestbook) List(ctx context.Context, req *apiv1.GuestbookListRequest) (*apiv1.GuestbookListResponse, error) {
	num := 10
	guestbooks, err := g.store.ListGuestbook(ctx, int(req.Page), num)
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
		}
		items = append(items, item)
	}

	total, err := g.store.CountGuestbookTotal(ctx)
	if err != nil {
		return nil, err
	}

	return &apiv1.GuestbookListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (g *Guestbook) Create(ctx context.Context, req *apiv1.GuestbookCreateRequest) (*apiv1.GuestbookCreateResponse, error) {
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
