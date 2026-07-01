package admin

import (
	"context"
	"time"

	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.LinkServiceHTTPServer = (*Link)(nil)

type Link struct {
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
		items = append(items, adminv1.LinkItem_builder{Id: int32(v.Id),
			Name:      v.Name,
			Url:       v.Url,
			Desc:      v.Desc,
			Status:    v.Status,
			CreatedAt: v.CreatedAt.Format(time.DateTime)}.Build(),
		)
	}
	return adminv1.LinkListResponse_builder{List: items,
			Total: int32(len(items))}.Build(),
		nil
}

func (l *Link) Create(ctx context.Context, req *adminv1.LinkCreateRequest) (*types.IDResponse, error) {
	m := &model.Link{
		Name:      req.GetName(),
		Url:       req.GetUrl(),
		Desc:      req.GetDesc(),
		Status:    model.LinkStatusPending,
		CreatedAt: time.Now(),
	}
	lastId, err := l.store.CreateLink(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(lastId)}.Build(), nil
}

func (l *Link) Update(ctx context.Context, req *adminv1.LinkUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateLink{Id: int(req.GetId())}
	if req.GetName() != "" {
		v := req.GetName()
		u.Name = &v
	}
	if req.GetUrl() != "" {
		v := req.GetUrl()
		u.Url = &v
	}
	if req.GetDesc() != "" {
		v := req.GetDesc()
		u.Desc = &v
	}
	if err := l.store.UpdateLink(ctx, u); err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: req.GetId()}.Build(), nil
}

func (l *Link) Delete(ctx context.Context, req *adminv1.LinkDeleteRequest) (*emptypb.Empty, error) {
	if err := l.store.DeleteLink(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Approve 审核链接（通过或驳回）
func (l *Link) Approve(ctx context.Context, req *adminv1.LinkApproveRequest) (*emptypb.Empty, error) {
	status := req.GetStatus()
	u := &model.UpdateLink{
		Id:     int(req.GetId()),
		Status: &status,
	}
	if err := l.store.UpdateLink(ctx, u); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
