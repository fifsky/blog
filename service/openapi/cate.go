package openapi

import (
	"context"
	"strconv"

	apiv1 "app/proto/gen/api/v1"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.CateServiceServer = (*Cate)(nil)

type Cate struct {
	apiv1.UnimplementedCateServiceServer
	store *store.Store
}

func NewCate(s *store.Store) *Cate {
	return &Cate{store: s}
}

func (c *Cate) All(ctx context.Context, _ *emptypb.Empty) (*apiv1.CateMenuResponse, error) {
	cates, err := c.store.GetAllCates(ctx)
	if err != nil {
		return nil, err
	}
	resp := apiv1.CateMenuResponse_builder{}.Build()
	for _, v := range cates {
		resp.SetList(append(resp.GetList(), apiv1.CateMenuItem_builder{Url: "/category/" + v.Domain,
			Content: v.Name + "(" + strconv.Itoa(v.Num) + ")"}.Build(),
		))
	}
	return resp, nil
}
