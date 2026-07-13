package admin

import (
	"context"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.GuestbookServiceHTTPServer = (*Guestbook)(nil)

type Guestbook struct {
	store *store.Store
}

func NewGuestbook(s *store.Store) *Guestbook {
	return &Guestbook{store: s}
}

func (g *Guestbook) List(ctx context.Context, req *adminv1.GuestbookListRequest) (*adminv1.GuestbookListResponse, error) {
	num := 10
	guestbooks, err := g.store.ListGuestbook(ctx, req.GetKeyword(), int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}

	items := lo.Map(guestbooks, func(gb model.Guestbook, _ int) *adminv1.GuestbookItem {
		return adminv1.GuestbookItem_builder{Id: int32(gb.Id),
			Name:      gb.Name,
			Content:   gb.Content,
			Ip:        gb.Ip,
			Top:       int32(gb.Top),
			CreatedAt: gb.CreatedAt.Format("2006-01-02 15:04:05")}.Build()
	})

	total, err := g.store.CountGuestbookTotal(ctx, req.GetKeyword())
	if err != nil {
		return nil, err
	}

	return adminv1.GuestbookListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (g *Guestbook) Delete(ctx context.Context, req *adminv1.GuestbookDeleteRequest) (*emptypb.Empty, error) {
	ids := lo.Map(req.GetIds(), func(id int32, _ int) int { return int(id) })
	if err := g.store.DeleteGuestbook(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
