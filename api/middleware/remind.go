package middleware

import (
	"app/provider/model"
	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/golib/convert"
	"github.com/ilibs/gosql/v2"

	"app/config"
	"app/pkg/aesutil"
)

type RemindAuth gee.HandlerFunc

func NewRemindAuth(db *gosql.DB, conf *config.Config) RemindAuth {
	return func(c *gee.Context) gee.Response {
		token := c.Query("token")

		if token == "" {
			c.Abort()
			return response.Fail(c, 201, "非法访问")
		}

		id, err := aesutil.AesDecode(conf.Common.TokenSecret, token)
		if err != nil {
			c.Abort()
			return response.Fail(c, 202, "Token错误")
		}

		remind := &model.Reminds{Id: convert.StrTo(id).MustInt()}
		err = db.Model(remind).Get()
		if err != nil {
			c.Abort()
			return response.Fail(c, 203, "数据不存在")
		}

		c.Set("remind", remind)
		c.Next()
		return nil
	}
}
