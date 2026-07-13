package admin

import (
	"context"
	"time"

	apperrors "app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.CateServiceHTTPServer = (*Cate)(nil)

type Cate struct {
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
	items := lo.Map(cates, func(v model.CateArtivleCount, _ int) *adminv1.CateListItem {
		return adminv1.CateListItem_builder{Id: int32(v.Id),
			Name:      v.Name,
			Desc:      v.Desc,
			Domain:    v.Domain,
			CreatedAt: v.CreatedAt.Format(time.DateTime),
			UpdatedAt: v.UpdatedAt.Format(time.DateTime),
			Num:       int32(v.Num)}.Build()
	})
	return adminv1.CateListResponse_builder{List: items,
			Total: int32(len(items))}.Build(),
		nil
}

func (c *Cate) Create(ctx context.Context, req *adminv1.CateCreateRequest) (*types.IDResponse, error) {
	now := time.Now()
	m := &model.Cate{
		Name:      req.GetName(),
		Desc:      req.GetDesc(),
		Domain:    req.GetDomain(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := c.store.CreateCate(ctx, m)
	if err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(lastId)}.Build(), nil
}

func (c *Cate) Update(ctx context.Context, req *adminv1.CateUpdateRequest) (*types.IDResponse, error) {
	now := time.Now()
	u := &model.UpdateCate{Id: int(req.GetId())}
	if req.GetName() != "" {
		u.Name = new(req.GetName())
	}
	if req.GetDesc() != "" {
		u.Desc = new(req.GetDesc())
	}
	if req.GetDomain() != "" {
		u.Domain = new(req.GetDomain())
	}
	u.UpdatedAt = &now
	if err := c.store.UpdateCate(ctx, u); err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(req.GetId())}.Build(), nil
}

func (c *Cate) Delete(ctx context.Context, req *adminv1.CateDeleteRequest) (*emptypb.Empty, error) {
	total, _ := c.store.PostsCount(ctx, int(req.GetId()))
	if total > 0 {
		return nil, apperrors.BadRequest("CATE_HAS_POSTS", "该分类下面还有文章，不能删除")
	}
	if err := c.store.DeleteCate(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
