package openapi

import (
	"context"
	"crypto/md5"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"app/config"
	"app/pkg/gotp"
	"app/pkg/ipgeo"
	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/service/feishu"
	"app/store"
	"app/store/model"

	"github.com/goapt/logger"
)

var _ apiv1.UserServiceHTTPServer = (*User)(nil)

type User struct {
	store      *store.Store
	conf       *config.Config
	notifyCard *feishu.NotifyCard
	sender     *feishu.FeishuSender
	httpClient *http.Client
}

func NewUser(s *store.Store, conf *config.Config, sender *feishu.FeishuSender, httpClient *http.Client) *User {
	return &User{
		store:      s,
		conf:       conf,
		notifyCard: feishu.NewNotifyCard(),
		sender:     sender,
		httpClient: httpClient,
	}
}

func (u *User) Login(ctx context.Context, in *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	user, err := u.store.GetUserByName(ctx, in.GetUserName())
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(in.GetPassword()))) {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if user.Status != model.UserStatusActive {
		return nil, fmt.Errorf("用户已停用")
	}
	if user.TotpSecret != "" {
		if in.GetTotpCode() == "" {
			return apiv1.LoginResponse_builder{RequireTotp: true}.Build(), nil
		}
		totp := gotp.NewDefaultTOTP(user.TotpSecret)
		ok, err := totp.Verify(in.GetTotpCode(), int64(time.Now().Unix()))
		if err != nil || !ok {
			return nil, fmt.Errorf("2FA验证码错误")
		}
	}
	tokenString, expiresAt, err := signAccessToken(u.conf.Common.TokenSecret, user.Id)
	if err != nil {
		return nil, fmt.Errorf("Access Token加密错误:%s", err)
	}

	// 异步发送登录通知
	u.notifyLogin(ctx, user.Name)

	return apiv1.LoginResponse_builder{AccessToken: tokenString,
			ExpiresAt: expiresAt,
			User: apiv1.UserItem_builder{Id: int32(user.Id),
				Name:      user.Name,
				NickName:  user.NickName,
				Email:     user.Email,
				Status:    string(user.Status),
				Type:      int32(user.Type),
				CreatedAt: user.CreatedAt.Format(time.DateTime),
				UpdatedAt: user.UpdatedAt.Format(time.DateTime)}.Build()}.Build(),

		nil
}

// notifyLogin 异步发送登录通知到飞书
func (u *User) notifyLogin(ctx context.Context, userName string) {
	if u.sender == nil || u.notifyCard == nil {
		return
	}

	ip := middleware.ClientIPFromContext(ctx)

	go func() {
		region := ""
		if ip != "" {
			geo, err := ipgeo.Lookup(context.Background(), u.httpClient, ip)
			if err != nil {
				logger.Error("login notify ipgeo lookup error", slog.String("err", err.Error()), slog.String("ip", ip))
			} else {
				if geo.City != "" {
					region = geo.City
				}
				if geo.Region != "" {
					if region != "" {
						region += ", " + geo.Region
					} else {
						region = geo.Region
					}
				}
				if geo.Country != "" {
					if region != "" {
						region += ", " + geo.Country
					} else {
						region = geo.Country
					}
				}
			}
		}

		content := fmt.Sprintf("%s 登录了博客", userName)
		if ip != "" {
			content += fmt.Sprintf("\n登录IP: %s", ip)
		}
		if region != "" {
			content += fmt.Sprintf("\n地区: %s", region)
		}

		msg := feishu.NotifyMessage{
			Content: content,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
		}
		cardJSON := u.notifyCard.BuildCard(msg)

		if err := u.sender.Send(context.Background(), cardJSON); err != nil {
			logger.Error("login notify send error", slog.String("err", err.Error()))
		}
	}()
}
