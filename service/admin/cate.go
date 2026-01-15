package admin

import (
	"context"
	"fmt"
	"time"

	"app/model"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.CateServiceServer = (*Cate)(nil)

type Cate struct {
	adminv1.UnimplementedCateServiceServer
	store *store.Store
}

func NewCate(s *store.Store) *Cate {
	return &Cate{store: s}
}

func (c *Cate) List(ctx context.Context, _ *emptypb.Empty) (*adminv1.CateListResponse, error) {
	cates, err := c.store.GetAllCates(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.CateListItem, 0, len(cates))
	for _, v := range cates {
		items = append(items, &adminv1.CateListItem{
			Id:        int32(v.Id),
			Name:      v.Name,
			Desc:      v.Desc,
			Domain:    v.Domain,
			CreatedAt: v.CreatedAt.Format(time.DateTime),
			UpdatedAt: v.UpdatedAt.Format(time.DateTime),
			Num:       int32(v.Num),
		})
	}
	return &adminv1.CateListResponse{
		List:  items,
		Total: int32(len(items)),
	}, nil
}

func (c *Cate) Create(ctx context.Context, req *adminv1.CateCreateRequest) (*types.IDResponse, error) {
	now := time.Now()
	m := &model.Cate{
		Name:      req.Name,
		Desc:      req.Desc,
		Domain:    req.Domain,
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := c.store.CreateCate(ctx, m)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastId)}, nil
}

func (c *Cate) Update(ctx context.Context, req *adminv1.CateUpdateRequest) (*types.IDResponse, error) {
	now := time.Now()
	u := &model.UpdateCate{Id: int(req.Id)}
	if req.Name != "" {
		v := req.Name
		u.Name = &v
	}
	if req.Desc != "" {
		v := req.Desc
		u.Desc = &v
	}
	if req.Domain != "" {
		v := req.Domain
		u.Domain = &v
	}
	u.UpdatedAt = &now
	if err := c.store.UpdateCate(ctx, u); err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(req.Id)}, nil
}

func (c *Cate) Delete(ctx context.Context, req *adminv1.CateDeleteRequest) (*emptypb.Empty, error) {
	total, _ := c.store.PostsCount(ctx, int(req.Id))
	if total > 0 {
		return nil, fmt.Errorf("该分类下面还有文章，不能删除")
	}
	if err := c.store.DeleteCate(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
