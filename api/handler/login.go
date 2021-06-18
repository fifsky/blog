package handler

import (
	"fmt"

	"app/config"
	"app/model"
	"app/pkg/aesutil"
	"app/response"

	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/hashing"
	"github.com/ilibs/gosql/v2"
)

var Login gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, err)
	}

	if p.UserName == "" || p.Password == "" {
		return response.Fail(c, 201, "用户名密码不能为空")
	}

	user := &model.Users{Name: p.UserName, Password: hashing.Md5(p.Password)}
	err := gosql.Model(user).Get()
	if err != nil {
		return response.Fail(c, 202, "用户名或密码错误")
	}

	if user.Status != 1 {
		return response.Fail(c, 202, "用户已停用")
	}

	src := fmt.Sprintf("%d:%s", user.Id, hashing.Md5(fmt.Sprintf("%d%s", user.Id, config.App.Common.TokenSecret)))
	cipherText, err := aesutil.AesEncode(config.App.Common.TokenSecret, src)
	if err != nil {
		return response.Fail(c, 201, "Access Token加密错误"+err.Error())
	}

	return response.Success(c, gin.H{
		"access_token": cipherText,
		"user":         user,
	})
}
