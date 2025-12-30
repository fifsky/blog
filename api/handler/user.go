package handler

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"app/config"
	"app/model"
	"app/response"
	"app/store"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	store *store.Store
	conf  *config.Config
}

func NewUser(s *store.Store, conf *config.Config) *User {
	return &User{store: s, conf: conf}
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	// 解析登录请求参数
	p, err := decode[LoginRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	user, err := u.store.GetUserByName(r.Context(), p.UserName)
	if err != nil {
		response.Fail(w, 202, "用户名或密码错误")
		return
	}

	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(p.Password))) {
		response.Fail(w, 202, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		response.Fail(w, 202, "用户已停用")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    fmt.Sprint(user.Id),
	})

	tokenString, err := token.SignedString([]byte(u.conf.Common.TokenSecret))

	if err != nil {
		response.Fail(w, 201, "Access Token加密错误"+err.Error())
		return
	}
	ui := UserItem{
		Id:        user.Id,
		Name:      user.Name,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    user.Status,
		Type:      user.Type,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	response.Success(w, map[string]any{
		"access_token": tokenString,
		"user":         ui,
	})
	return
}

func (u *User) LoginUser(w http.ResponseWriter, r *http.Request) {
	// 从请求上下文中获取登录用户
	user := getLoginUser(r.Context())
	ui := UserItem{
		Id:        user.Id,
		Name:      user.Name,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    user.Status,
		Type:      user.Type,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	response.Success(w, ui)
}

func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	// 解析ID参数
	p, err := decode[IDRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	user, err := u.store.GetUser(r.Context(), p.Id)
	if err != nil {
		response.Fail(w, 201, "用户不存在")
		return
	}

	ui := UserItem{
		Id:        user.Id,
		Name:      user.Name,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    user.Status,
		Type:      user.Type,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	response.Success(w, ui)
}

func (u *User) List(w http.ResponseWriter, r *http.Request) {
	// 解析分页参数
	p, err := decode[PageRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	num := 10
	users, err := u.store.ListUser(r.Context(), p.Page, num)
	if err != nil {
		response.Fail(w, 202, err)
		return
	}
	items := make([]UserItem, 0, len(users))
	for _, user := range users {
		items = append(items, UserItem{
			Id:        user.Id,
			Name:      user.Name,
			NickName:  user.NickName,
			Email:     user.Email,
			Status:    user.Status,
			Type:      user.Type,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	total, err := u.store.CountUserTotal(r.Context())
	if err != nil {
		response.Fail(w, 500, err)
		return
	}

	resp := UserListReponse{
		List:      items,
		PageTotal: totalPages(total, num),
	}

	response.Success(w, resp)
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	in, err := decode[UserCreateRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	if in.Password == "" {
		response.Fail(w, 201, "密码不能为空")
		return
	}
	hashed := fmt.Sprintf("%x", md5.Sum([]byte(in.Password)))
	now := time.Now()
	cReq := &model.User{
		Name:      in.Name,
		Password:  hashed,
		NickName:  in.NickName,
		Email:     in.Email,
		Status:    1,
		Type:      in.Type,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := u.store.CreateUser(r.Context(), cReq); err != nil {
		response.Fail(w, 201, "创建失败")
		return
	}
	resp := UserPostResponse{
		Id:       0,
		Name:     in.Name,
		NickName: in.NickName,
		Email:    in.Email,
		Status:   1,
		Type:     in.Type,
	}
	response.Success(w, resp)
}

func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	in, err := decode[UserUpdateRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	if in.Id <= 0 {
		response.Fail(w, 201, "参数错误: ID不能为空")
		return
	}
	hashed := fmt.Sprintf("%x", md5.Sum([]byte(in.Password)))
	now := time.Now()
	uReq := &model.UpdateUser{
		Id:        in.Id,
		UpdatedAt: &now,
	}
	if in.Name != "" {
		uReq.Name = &in.Name
	}
	if in.Password != "" {
		uReq.Password = &hashed
	}
	if in.NickName != "" {
		uReq.NickName = &in.NickName
	}
	if in.Email != "" {
		uReq.Email = &in.Email
	}
	if in.Type > 0 {
		uReq.Type = &in.Type
	}
	if err := u.store.UpdateUser(r.Context(), uReq); err != nil {
		response.Fail(w, 201, "更新失败")
		return
	}
	resp := UserPostResponse{
		Id:       in.Id,
		Name:     in.Name,
		NickName: in.NickName,
		Email:    in.Email,
		Status:   1,
		Type:     in.Type,
	}
	response.Success(w, resp)
}

func (u *User) Status(w http.ResponseWriter, r *http.Request) {
	// 解析ID参数
	p, err := decode[IDRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	user, err := u.store.GetUser(r.Context(), p.Id)
	if err != nil {
		response.Fail(w, 202, "用户不存在")
		return
	}

	status := user.Status
	if status == 1 {
		status = 2
	} else {
		status = 1
	}

	if err := u.store.UpdateUser(r.Context(), &model.UpdateUser{Id: p.Id, Status: &status}); err != nil {
		response.Fail(w, 201, "停启用失败")
		return
	}
	response.Success(w, nil)
}
