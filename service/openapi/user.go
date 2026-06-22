package openapi

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"app/config"
	"app/pkg/gotp"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
)

var _ apiv1.UserServiceHTTPServer = (*User)(nil)

type User struct {
	store *store.Store
	conf  *config.Config
}

func NewUser(s *store.Store, conf *config.Config) *User {
	return &User{
		store: s,
		conf:  conf,
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
	if user.Status != 1 {
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
	return apiv1.LoginResponse_builder{AccessToken: tokenString,
			ExpiresAt: expiresAt,
			User: apiv1.UserItem_builder{Id: int32(user.Id),
				Name:      user.Name,
				NickName:  user.NickName,
				Email:     user.Email,
				Status:    int32(user.Status),
				Type:      int32(user.Type),
				CreatedAt: user.CreatedAt.Format(time.DateTime),
				UpdatedAt: user.UpdatedAt.Format(time.DateTime)}.Build()}.Build(),

		nil
}
