package handler

import (
	"app/response"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/hashing"
	"github.com/goapt/golib/pagination"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"

	"app/model"
)

var AdminLoginUser gee.HandlerFunc = func(c *gee.Context) gee.Response {
	user := getLoginUser(c)
	return response.Success(c, user)
}

var AdminUserGet gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	user := &model.Users{Id: p.Id}
	err := gosql.Model(user).Get()
	if err != nil {
		return response.Fail(c, 201, "用户不存在")
	}

	return response.Success(c, user)
}

var AdminUserList gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	h := gin.H{}
	num := 10
	users, err := model.UserGetList(p.Page, num)
	h["list"] = users

	total, err := gosql.Model(&model.Users{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, h)
}

var AdminUserPost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	users := &model.Users{}
	if err := c.ShouldBindJSON(users); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if users.Name == "" {
		return response.Fail(c, 201, "用户名不能为空")
	}

	if users.Id == 0 && users.Password == "" {
		return response.Fail(c, 201, "密码不能为空")
	} else {
		users.Password = hashing.Md5(users.Password)
	}

	if users.NickName == "" {
		return response.Fail(c, 201, "昵称不能为空")
	}

	if users.Id > 0 {
		if _, err := gosql.Model(users).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(users).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "创建失败")
		}
	}

	return response.Success(c, users)
}

var AdminUserStatus gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}
	user := &model.Users{Id: p.Id}
	err := gosql.Model(user).Get()
	if err != nil {
		return response.Fail(c, 202, "用户不存在:"+err.Error())
	}

	status := user.Status
	if status == 1 {
		status = 2
	} else {
		status = 1
	}

	if _, err := gosql.Model(&model.Users{Status: status}).Where("id = ?", p.Id).Update(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "停启用失败")
	}
	return response.Success(c, nil)
}
