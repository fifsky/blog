package admin

import (
	"context"

	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"
)

var _ adminv1.SettingServiceServer = (*Setting)(nil)

type Setting struct {
	adminv1.UnimplementedSettingServiceServer
	store *store.Store
}

func NewSetting(s *store.Store) *Setting {
	return &Setting{store: s}
}

func (s *Setting) Update(ctx context.Context, req *types.Options) (*types.Options, error) {
	m, err := s.store.UpdateOptions(ctx, req.Kv)
	if err != nil {
		return nil, err
	}
	return &types.Options{Kv: m}, nil
}
