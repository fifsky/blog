package openapi

import (
	"context"
	"time"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"app/store/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.LinkServiceHTTPServer = (*Link)(nil)

type Link struct {
	store *store.Store
}

func NewLink(s *store.Store) *Link {
	return &Link{store: s}
}

// All 获取所有已审核通过的链接
func (l *Link) All(ctx context.Context, _ *emptypb.Empty) (*apiv1.LinkMenuResponse, error) {
	links, err := l.store.GetApprovedLinks(ctx)
	if err != nil {
		return nil, err
	}
	resp := apiv1.LinkMenuResponse_builder{}.Build()
	for _, v := range links {
		resp.SetList(append(resp.GetList(), apiv1.LinkMenuItem_builder{Url: v.Url,
			Content: v.Name}.Build(),
		))
	}
	return resp, nil
}

// Submit 提交友情链接申请（状态为审核中）
func (l *Link) Submit(ctx context.Context, req *apiv1.LinkSubmitRequest) (*emptypb.Empty, error) {
	m := &model.Link{
		Name:      req.GetName(),
		Url:       req.GetUrl(),
		Desc:      req.GetDesc(),
		Status:    model.LinkStatusPending,
		CreatedAt: time.Now(),
	}
	_, err := l.store.CreateLink(ctx, m)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
