package admin

import (
	"context"
	"time"

	"app/model"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.LinkServiceServer = (*Link)(nil)

type Link struct {
	adminv1.UnimplementedLinkServiceServer
	store *store.Store
}

func NewLink(s *store.Store) *Link {
	return &Link{store: s}
}

func (l *Link) List(ctx context.Context, _ *emptypb.Empty) (*adminv1.LinkListResponse, error) {
	links, err := l.store.GetAllLinks(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.LinkItem, 0, len(links))
	for _, v := range links {
		items = append(items, &adminv1.LinkItem{
			Id:        int32(v.Id),
			Name:      v.Name,
			Url:       v.Url,
			Desc:      v.Desc,
			CreatedAt: v.CreatedAt.Format(time.DateTime),
		})
	}
	return &adminv1.LinkListResponse{
		List:  items,
		Total: int32(len(items)),
	}, nil
}

func (l *Link) Create(ctx context.Context, req *adminv1.LinkCreateRequest) (*adminv1.IDResponse, error) {
	m := &model.Link{
		Name:      req.Name,
		Url:       req.Url,
		Desc:      req.Desc,
		CreatedAt: time.Now(),
	}
	lastId, err := l.store.CreateLink(ctx, m)
	if err != nil {
		return nil, err
	}
	return &adminv1.IDResponse{Id: int32(lastId)}, nil
}

func (l *Link) Update(ctx context.Context, req *adminv1.LinkUpdateRequest) (*adminv1.IDResponse, error) {
	u := &model.UpdateLink{Id: int(req.Id)}
	if req.Name != "" {
		v := req.Name
		u.Name = &v
	}
	if req.Url != "" {
		v := req.Url
		u.Url = &v
	}
	if req.Desc != "" {
		v := req.Desc
		u.Desc = &v
	}
	if err := l.store.UpdateLink(ctx, u); err != nil {
		return nil, err
	}
	return &adminv1.IDResponse{Id: req.Id}, nil
}

func (l *Link) Delete(ctx context.Context, req *adminv1.IDRequest) (*emptypb.Empty, error) {
	if err := l.store.DeleteLink(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
