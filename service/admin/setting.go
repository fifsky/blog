package admin

import (
	"context"
	"strconv"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.SettingServiceServer = (*Setting)(nil)

type Setting struct {
	adminv1.UnimplementedSettingServiceServer
	store *store.Store
}

func NewSetting(s *store.Store) *Setting {
	return &Setting{store: s}
}

func (s *Setting) Get(ctx context.Context, _ *emptypb.Empty) (*adminv1.AdminSetting, error) {
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

	return &adminv1.AdminSetting{
		SiteName:    m["site_name"],
		SiteDesc:    m["site_desc"],
		SiteKeyword: m["site_keyword"],
		PostNum:     int32(postNum),
		AiEndpoint:  m["ai_endpoint"],
		AiModel:     m["ai_model"],
		AiToken:     m["ai_token"],
	}, nil
}

func (s *Setting) Update(ctx context.Context, req *adminv1.AdminSetting) (*adminv1.AdminSetting, error) {
	kv := map[string]string{
		"site_name":    req.SiteName,
		"site_desc":    req.SiteDesc,
		"site_keyword": req.SiteKeyword,
		"post_num":     strconv.Itoa(int(req.PostNum)),
		"ai_endpoint":  req.AiEndpoint,
		"ai_model":     req.AiModel,
		"ai_token":     req.AiToken,
	}

	m, err := s.store.UpdateOptions(ctx, kv)
	if err != nil {
		return nil, err
	}
	
	postNum := 10
	if val, ok := m["post_num"]; ok {
		if n, err := strconv.Atoi(val); err == nil {
			postNum = n
		}
	}

	return &adminv1.AdminSetting{
		SiteName:    m["site_name"],
		SiteDesc:    m["site_desc"],
		SiteKeyword: m["site_keyword"],
		PostNum:     int32(postNum),
		AiEndpoint:  m["ai_endpoint"],
		AiModel:     m["ai_model"],
		AiToken:     m["ai_token"],
	}, nil
}
