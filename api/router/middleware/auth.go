package middleware

import (
	"strings"

	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/logger"

	"app/config"
	"app/model"
	"app/pkg/aesutil"

	"github.com/goapt/golib/hashing"
	"github.com/ilibs/gosql/v2"
)

var AuthLogin gee.HandlerFunc = func(c *gee.Context) gee.Response {
	accessToken := c.Request.Header.Get("Access-Token")

	if accessToken == "" {
		c.Abort()
		return response.Fail(c, 201, "Access Token不能为空")
	}

	cipherText, err := aesutil.AesDecode(config.App.Common.TokenSecret, accessToken)
	if err != nil {
		logger.Data(map[string]interface{}{
			"token": accessToken,
			"err":   err,
		}).Error("Access Token错误")
		c.Abort()
		return response.Fail(c, 201, "Access Token错误")
	}

	v := strings.Split(cipherText, ":")
	if len(v) != 2 || hashing.Md5(v[0]+config.App.Common.TokenSecret) != v[1] {
		logger.Data(map[string]interface{}{
			"token":      accessToken,
			"cipherText": cipherText,
		}).Error("Access Token不合法")
		c.Abort()
		return response.Fail(c, 201, "Access Token不合法")
	}

	user := &model.Users{}
	err = gosql.Model(user).Where("id = ?", v[0]).Get()
	if err != nil {
		c.Abort()
		return response.Fail(c, 201, "Access Token错误，用户不存在")
	}

	c.Set("userInfo", user)
	c.Next()
	return nil
}
