package service

import (
	"context"
	"time"

	"app/config"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
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

func (u *User) Get(ctx context.Context, request *apiv1.GetUserRequest) (*apiv1.User, error) {
	user, err := u.store.GetUser(ctx, int(request.Id))
	if err != nil {
		return nil, err
	}

	return &apiv1.User{
		Id:        int32(user.Id),
		Name:      user.Name,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    int32(user.Status),
		Type:      int32(user.Type),
		CreatedAt: user.CreatedAt.Format(time.DateTime),
		UpdatedAt: user.UpdatedAt.Format(time.DateTime),
	}, nil
}
