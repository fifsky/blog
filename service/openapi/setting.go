package openapi

import (
	"context"
	"io"
	"net/http"

	apiv1 "app/proto/gen/api/v1"
	"app/proto/gen/types"
	"app/store"

	"google.golang.org/genproto/googleapis/api/httpbody"
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

func (s *Setting) GetChinaMap(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	resp, err := http.Get("https://geo.datav.aliyun.com/areas_v3/bound/100000_full.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        body,
	}, nil
}
