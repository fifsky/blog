package admin

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"app/config"
	"app/model"
	adminv1 "app/proto/gen/admin/v1"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.UserServiceServer = (*User)(nil)

type User struct {
	adminv1.UnimplementedUserServiceServer
	store *store.Store
	conf  *config.Config
}

func NewUser(s *store.Store, conf *config.Config) *User {
	return &User{
		store: s,
		conf:  conf,
	}
}

func (u *User) Get(ctx context.Context, request *adminv1.GetUserRequest) (*adminv1.User, error) {
	user, err := u.store.GetUser(ctx, int(request.Id))
	if err != nil {
		return nil, err
	}

	return &adminv1.User{
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

func (u *User) Create(ctx context.Context, in *adminv1.UserCreateRequest) (*adminv1.IDResponse, error) {
	if in.Password == "" {
		return nil, fmt.Errorf("密码不能为空")
	}
	hashed := fmt.Sprintf("%x", md5.Sum([]byte(in.Password)))
	now := time.Now()
	cReq := &model.User{
		Name:      in.Name,
		Password:  hashed,
		NickName:  in.NickName,
		Email:     in.Email,
		Status:    1,
		Type:      int(in.Type),
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := u.store.CreateUser(ctx, cReq)
	if err != nil {
		return nil, err
	}
	return &adminv1.IDResponse{Id: int32(lastId)}, nil
}

func (u *User) Update(ctx context.Context, in *adminv1.UserUpdateRequest) (*adminv1.IDResponse, error) {
	hashed := fmt.Sprintf("%x", md5.Sum([]byte(in.Password)))
	now := time.Now()
	uReq := &model.UpdateUser{
		Id:        int(in.Id),
		UpdatedAt: &now,
	}
	if in.Name != "" {
		v := in.Name
		uReq.Name = &v
	}
	if in.Password != "" {
		v := hashed
		uReq.Password = &v
	}
	if in.NickName != "" {
		v := in.NickName
		uReq.NickName = &v
	}
	if in.Email != "" {
		v := in.Email
		uReq.Email = &v
	}
	if in.Type > 0 {
		v := int(in.Type)
		uReq.Type = &v
	}
	if err := u.store.UpdateUser(ctx, uReq); err != nil {
		return nil, err
	}
	return &adminv1.IDResponse{Id: int32(in.Id)}, nil
}

func (u *User) List(ctx context.Context, req *adminv1.PageRequest) (*adminv1.UserListResponse, error) {
	num := 10
	users, err := u.store.ListUser(ctx, int(req.Page), num)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.UserItem, 0, len(users))
	for _, user := range users {
		items = append(items, &adminv1.UserItem{
			Id:        int32(user.Id),
			Name:      user.Name,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    int32(user.Status),
			Type:      int32(user.Type),
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			UpdatedAt: user.UpdatedAt.Format(time.DateTime),
		})
	}
	total, err := u.store.CountUserTotal(ctx)
	if err != nil {
		return nil, err
	}
	return &adminv1.UserListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (u *User) Status(ctx context.Context, req *adminv1.IDRequest) (*emptypb.Empty, error) {
	user, err := u.store.GetUser(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	status := user.Status
	if status == 1 {
		status = 2
	} else {
		status = 1
	}
	if err := u.store.UpdateUser(ctx, &model.UpdateUser{Id: int(req.Id), Status: &status}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (u *User) LoginUser(ctx context.Context, _ *emptypb.Empty) (*adminv1.User, error) {
	user := GetLoginUser(ctx)
	return &adminv1.User{
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
