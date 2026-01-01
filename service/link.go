package service

import (
	"context"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.LinkServiceServer = (*Link)(nil)

type Link struct {
	apiv1.UnimplementedLinkServiceServer
	store *store.Store
}

func NewLink(s *store.Store) *Link {
	return &Link{store: s}
}

func (l *Link) All(ctx context.Context, _ *emptypb.Empty) (*apiv1.LinkMenuResponse, error) {
	links, err := l.store.GetAllLinks(ctx)
	if err != nil {
		return nil, err
	}
	resp := &apiv1.LinkMenuResponse{}
	for _, v := range links {
		resp.List = append(resp.List, &apiv1.LinkMenuItem{
			Url:     v.Url,
			Content: v.Name,
		})
	}
	return resp, nil
}
