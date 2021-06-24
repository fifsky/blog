package middleware

import (
	"fmt"

	model2 "app/provider/model"
	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/golib/convert"
	"github.com/ilibs/gosql/v2"

	"app/config"
	"app/pkg/aesutil"
)

type RemindAuth gee.HandlerFunc

func NewRemindAuth() RemindAuth {
	return func(c *gee.Context) gee.Response {
		token := c.Query("token")

		if token == "" {
			c.Abort()
			return response.Fail(c, 201, "非法访问")
		}

		id, err := aesutil.AesDecode(config.App.Common.TokenSecret, token)
		if err != nil {
			c.Abort()
			return response.Fail(c, 202, fmt.Sprintf("error:%s token:%s", err, token))
		}

		remind := &model2.Reminds{Id: convert.StrTo(id).MustInt()}
		err = gosql.Model(remind).Get()
		if err != nil {
			c.Abort()
			return response.Fail(c, 203, fmt.Sprintf("error:%s token:%s", err, token))
		}

		c.Set("remind", remind)
		c.Next()
		return nil
	}
}
