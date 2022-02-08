package handler

import (
	"fmt"
	"time"

	"app/config"
	"app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/golib/hashing"
	"github.com/goapt/golib/pagination"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilibs/gosql/v2"
)

type User struct {
	db       *gosql.DB
	userRepo *repo.User
	conf     *config.Config
}

func NewUser(db *gosql.DB, userRepo *repo.User, conf *config.Config) *User {
	return &User{db: db, userRepo: userRepo, conf: conf}
}

func (u *User) Login(c *gee.Context) gee.Response {
	p := &struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	user := &model.Users{Name: p.UserName, Password: hashing.Md5(p.Password)}
	err := u.db.Model(user).Get()
	if err != nil {
		return response.Fail(c, 202, "用户名或密码错误")
	}

	if user.Status != 1 {
		return response.Fail(c, 202, "用户已停用")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    fmt.Sprint(user.Id),
	})

	tokenString, err := token.SignedString([]byte(u.conf.Common.TokenSecret))

	if err != nil {
		return response.Fail(c, 201, "Access Token加密错误"+err.Error())
	}

	return response.Success(c, gee.H{
		"access_token": tokenString,
		"user":         user,
	})
}

func (u *User) LoginUser(c *gee.Context) gee.Response {
	user := getLoginUser(c)
	return response.Success(c, user)
}

func (u *User) Get(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	user := &model.Users{Id: p.Id}
	err := u.db.Model(user).Get()
	if err != nil {
		return response.Fail(c, 201, "用户不存在")
	}

	return response.Success(c, user)
}

func (u *User) List(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	h := gee.H{}
	num := 10
	users, err := u.userRepo.GetList(p.Page, num)
	if err != nil {
		return response.Fail(c, 202, err)
	}
	h["list"] = users

	total, err := u.db.Model(&model.Users{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, h)
}

func (u *User) Post(c *gee.Context) gee.Response {
	users := &model.Users{}
	if err := c.ShouldBindJSON(users); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if users.Id == 0 && users.Password == "" {
		return response.Fail(c, 201, "密码不能为空")
	} else {
		users.Password = hashing.Md5(users.Password)
	}

	if users.Id > 0 {
		if _, err := u.db.Model(users).Update(); err != nil {
			return response.Fail(c, 201, "更新失败")
		}
	} else {
		if _, err := u.db.Model(users).Create(); err != nil {
			return response.Fail(c, 201, "创建失败")
		}
	}

	return response.Success(c, users)
}

func (u *User) Status(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}
	user := &model.Users{Id: p.Id}
	err := u.db.Model(user).Get()
	if err != nil {
		return response.Fail(c, 202, "用户不存在")
	}

	status := user.Status
	if status == 1 {
		status = 2
	} else {
		status = 1
	}

	if _, err := u.db.Model(&model.Users{Status: status}).Where("id = ?", p.Id).Update(); err != nil {
		return response.Fail(c, 201, "停启用失败")
	}
	return response.Success(c, nil)
}
