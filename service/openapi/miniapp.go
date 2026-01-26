package openapi

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"app/config"
	"app/pkg/errors"
	"app/pkg/miniapp"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
)

var _ apiv1.MiniAppServiceServer = (*MiniApp)(nil)

type MiniApp struct {
	apiv1.UnimplementedMiniAppServiceServer
	store  *store.Store
	conf   *config.Config
	client *miniapp.Client
}

func NewMiniApp(s *store.Store, conf *config.Config, httpClient *http.Client) *MiniApp {
	return &MiniApp{
		store:  s,
		conf:   conf,
		client: miniapp.NewClient(conf.MiniAPP.Appid, conf.MiniAPP.AppSecret, httpClient),
	}
}

func (m *MiniApp) LoginCode(ctx context.Context, req *apiv1.MiniAppLoginRequest) (*apiv1.MiniAppLoginResponse, error) {
	if m.conf.MiniAPP.Appid == "" || m.conf.MiniAPP.AppSecret == "" {
		return nil, errors.ErrSystem.WithCause(fmt.Errorf("miniapp config missing"))
	}

	sess, err := m.client.Code2Session(ctx, req.Code)
	if err != nil {
		return nil, errors.BadRequest("MINIAPP_LOGIN_FAILED", "小程序登录失败").WithCause(err)
	}

	uid, err := m.store.GetUserIDByOpenid(ctx, sess.OpenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.BadRequest("MINIAPP_USER_NOT_FOUND", "该 openid 未绑定用户").WithMetadata(map[string]string{"openid": sess.OpenID})
		}
		return nil, errors.ErrSystem.WithCause(err)
	}

	user, err := m.store.GetUser(ctx, uid)
	if err != nil {
		return nil, errors.ErrSystem.WithCause(err)
	}
	if user.Status != 1 {
		return nil, errors.BadRequest("USER_DISABLED", "用户已停用")
	}

	token, err := signAccessToken(m.conf.Common.TokenSecret, user.Id)
	if err != nil {
		return nil, errors.ErrSystem.WithCause(err)
	}

	return &apiv1.MiniAppLoginResponse{
		AccessToken: token,
	}, nil
}
