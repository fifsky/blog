package openapi

import (
	"context"
	"io"
	"net/http"
	"strconv"

	apiv1 "app/proto/gen/api/v1"
	"app/store"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.SettingServiceHTTPServer = (*Setting)(nil)

type Setting struct {
	store *store.Store
}

func NewSetting(s *store.Store) *Setting {
	return &Setting{store: s}
}

func (s *Setting) Get(ctx context.Context, _ *emptypb.Empty) (*apiv1.Setting, error) {
	m, err := s.store.GetOptions(ctx)
	if err != nil {
		return nil, err
	}

	postNum := 10
	if val, ok := m["post_num"]; ok {
		if n, err := strconv.Atoi(val); err == nil {
			postNum = n
		}
	}

	return apiv1.Setting_builder{SiteName: m["site_name"],
			SiteDesc:    m["site_desc"],
			SiteKeyword: m["site_keyword"],
			PostNum:     int32(postNum)}.Build(),
		nil
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
