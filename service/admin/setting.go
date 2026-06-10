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

	return adminv1.AdminSetting_builder{SiteName: m["site_name"],
			SiteDesc:    m["site_desc"],
			SiteKeyword: m["site_keyword"],
			PostNum:     int32(postNum),
			AiEndpoint:  m["ai_endpoint"],
			AiModel:     m["ai_model"],
			AiToken:     m["ai_token"]}.Build(),
		nil
}

func (s *Setting) Update(ctx context.Context, req *adminv1.AdminSetting) (*adminv1.AdminSetting, error) {
	kv := map[string]string{
		"site_name":    req.GetSiteName(),
		"site_desc":    req.GetSiteDesc(),
		"site_keyword": req.GetSiteKeyword(),
		"post_num":     strconv.Itoa(int(req.GetPostNum())),
		"ai_endpoint":  req.GetAiEndpoint(),
		"ai_model":     req.GetAiModel(),
		"ai_token":     req.GetAiToken(),
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

	return adminv1.AdminSetting_builder{SiteName: m["site_name"],
			SiteDesc:    m["site_desc"],
			SiteKeyword: m["site_keyword"],
			PostNum:     int32(postNum),
			AiEndpoint:  m["ai_endpoint"],
			AiModel:     m["ai_model"],
			AiToken:     m["ai_token"]}.Build(),
		nil
}
