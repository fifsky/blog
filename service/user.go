package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"app/config"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
	"github.com/golang-jwt/jwt/v5"
)

var _ apiv1.UserServiceServer = (*User)(nil)

type User struct {
	apiv1.UnimplementedUserServiceServer
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
	user, err := u.store.GetUserByName(ctx, in.UserName)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(in.Password))) {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if user.Status != 1 {
		return nil, fmt.Errorf("用户已停用")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    fmt.Sprint(user.Id),
	})
	tokenString, err := token.SignedString([]byte(u.conf.Common.TokenSecret))
	if err != nil {
		return nil, fmt.Errorf("Access Token加密错误:%s", err)
	}
	return &apiv1.LoginResponse{
		AccessToken: tokenString,
		User: &apiv1.UserItem{
			Id:        int32(user.Id),
			Name:      user.Name,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    int32(user.Status),
			Type:      int32(user.Type),
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			UpdatedAt: user.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}
