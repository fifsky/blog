package admin

import (
	"context"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.GuestbookServiceServer = (*Guestbook)(nil)

type Guestbook struct {
	adminv1.UnimplementedGuestbookServiceServer
	store *store.Store
}

func NewGuestbook(s *store.Store) *Guestbook {
	return &Guestbook{store: s}
}

func (g *Guestbook) List(ctx context.Context, req *adminv1.GuestbookListRequest) (*adminv1.GuestbookListResponse, error) {
	num := 10
	guestbooks, err := g.store.ListGuestbook(ctx, req.Keyword, int(req.Page), num)
	if err != nil {
		return nil, err
	}

	items := make([]*adminv1.GuestbookItem, 0, len(guestbooks))
	for _, gb := range guestbooks {
		items = append(items, &adminv1.GuestbookItem{
			Id:        int32(gb.Id),
			Name:      gb.Name,
			Content:   gb.Content,
			Ip:        gb.Ip,
			Top:       int32(gb.Top),
			CreatedAt: gb.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	total, err := g.store.CountGuestbookTotal(ctx, req.Keyword)
	if err != nil {
		return nil, err
	}

	return &adminv1.GuestbookListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (g *Guestbook) Delete(ctx context.Context, req *adminv1.GuestbookDeleteRequest) (*emptypb.Empty, error) {
	ids := make([]int, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = int(id)
	}
	if err := g.store.DeleteGuestbook(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
