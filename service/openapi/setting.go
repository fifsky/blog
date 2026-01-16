package openapi

import (
	"context"

	apiv1 "app/proto/gen/api/v1"
	"app/proto/gen/types"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.SettingServiceServer = (*Setting)(nil)

type Setting struct {
	apiv1.UnimplementedSettingServiceServer
	store *store.Store
}

func NewSetting(s *store.Store) *Setting {
	return &Setting{store: s}
}

func (s *Setting) Get(ctx context.Context, _ *emptypb.Empty) (*types.Options, error) {
	m, err := s.store.GetOptions(ctx)
	if err != nil {
		return nil, err
	}
	return &types.Options{Kv: m}, nil
}
