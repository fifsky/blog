package admin

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"app/config"
	"app/pkg/gotp"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.UserServiceHTTPServer = (*User)(nil)

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

func (u *User) Get(ctx context.Context, request *adminv1.GetUserRequest) (*adminv1.User, error) {
	user, err := u.store.GetUser(ctx, int(request.GetId()))
	if err != nil {
		return nil, err
	}

	return adminv1.User_builder{Id: int32(user.Id),
			Name:      user.Name,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    int32(user.Status),
			Type:      int32(user.Type),
			HasTotp:   user.TotpSecret != "",
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			UpdatedAt: user.UpdatedAt.Format(time.DateTime)}.Build(),
		nil
}

func (u *User) Create(ctx context.Context, in *adminv1.UserCreateRequest) (*types.IDResponse, error) {
	if in.GetPassword() == "" {
		return nil, fmt.Errorf("密码不能为空")
	}
	hashed := fmt.Sprintf("%x", md5.Sum([]byte(in.GetPassword())))
	now := time.Now()
	cReq := &model.User{
		Name:      in.GetName(),
		Password:  hashed,
		NickName:  in.GetNickName(),
		Email:     in.GetEmail(),
		Status:    1,
		Type:      int(in.GetType()),
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := u.store.CreateUser(ctx, cReq)
	if err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(lastId)}.Build(), nil
}

func (u *User) Update(ctx context.Context, in *adminv1.UserUpdateRequest) (*types.IDResponse, error) {
	hashed := fmt.Sprintf("%x", md5.Sum([]byte(in.GetPassword())))
	now := time.Now()
	uReq := &model.UpdateUser{
		Id:        int(in.GetId()),
		UpdatedAt: &now,
	}
	if in.GetName() != "" {
		v := in.GetName()
		uReq.Name = &v
	}
	if in.GetPassword() != "" {
		v := hashed
		uReq.Password = &v
	}
	if in.GetNickName() != "" {
		v := in.GetNickName()
		uReq.NickName = &v
	}
	if in.GetEmail() != "" {
		v := in.GetEmail()
		uReq.Email = &v
	}
	if in.GetType() > 0 {
		v := int(in.GetType())
		uReq.Type = &v
	}
	if err := u.store.UpdateUser(ctx, uReq); err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(in.GetId())}.Build(), nil
}

func (u *User) List(ctx context.Context, req *adminv1.UserListRequest) (*adminv1.UserListResponse, error) {
	num := 10
	users, err := u.store.ListUser(ctx, int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.UserItem, 0, len(users))
	for _, user := range users {
		items = append(items, adminv1.UserItem_builder{Id: int32(user.Id),
			Name:      user.Name,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    int32(user.Status),
			Type:      int32(user.Type),
			HasTotp:   user.TotpSecret != "",
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			UpdatedAt: user.UpdatedAt.Format(time.DateTime)}.Build(),
		)
	}
	total, err := u.store.CountUserTotal(ctx)
	if err != nil {
		return nil, err
	}
	return adminv1.UserListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (u *User) Status(ctx context.Context, req *adminv1.UserStatusRequest) (*emptypb.Empty, error) {
	user, err := u.store.GetUser(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	status := user.Status
	if status == 1 {
		status = 2
	} else {
		status = 1
	}
	if err := u.store.UpdateUser(ctx, &model.UpdateUser{Id: int(req.GetId()), Status: &status}); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (u *User) LoginUser(ctx context.Context, _ *emptypb.Empty) (*adminv1.User, error) {
	user := GetLoginUser(ctx)
	return adminv1.User_builder{Id: int32(user.Id),
			Name:      user.Name,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    int32(user.Status),
			Type:      int32(user.Type),
			HasTotp:   user.TotpSecret != "",
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			UpdatedAt: user.UpdatedAt.Format(time.DateTime)}.Build(),
		nil
}

func (u *User) Generate2FA(ctx context.Context, req *adminv1.Generate2FARequest) (*adminv1.Generate2FAResponse, error) {
	user, err := u.store.GetUser(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	secret, err := gotp.RandomSecret(16)
	if err != nil {
		return nil, err
	}
	totp := gotp.NewDefaultTOTP(secret)
	issuer := "FIFSKY Blog"
	if u.conf.Env == "dev" || u.conf.Env == "local" || u.conf.Env == "development" {
		issuer = "FIFSKY Blog Dev"
	}
	uri, err := totp.ProvisioningUri(user.Name, issuer)
	if err != nil {
		return nil, err
	}
	return adminv1.Generate2FAResponse_builder{Secret: secret,
			QrCodeUri: uri}.Build(),
		nil
}

func (u *User) Bind2FA(ctx context.Context, req *adminv1.Bind2FARequest) (*emptypb.Empty, error) {
	if req.GetSecret() == "" || req.GetCode() == "" {
		return nil, fmt.Errorf("无效的绑定请求")
	}
	totp := gotp.NewDefaultTOTP(req.GetSecret())
	ok, err := totp.Verify(req.GetCode(), int64(time.Now().Unix()))
	if err != nil || !ok {
		return nil, fmt.Errorf("2FA验证码错误")
	}
	secret := req.GetSecret()
	err = u.store.UpdateUser(ctx, &model.UpdateUser{
		Id:         int(req.GetId()),
		TotpSecret: &secret,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
