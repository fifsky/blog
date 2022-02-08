package middleware

import (
	"app/provider/model"
	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/logger"

	"app/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ilibs/gosql/v2"
)

type AuthLogin gee.HandlerFunc

func NewAuthLogin(db *gosql.DB, conf *config.Config) AuthLogin {
	return func(c *gee.Context) gee.Response {
		accessToken := c.Request.Header.Get("Access-Token")

		if accessToken == "" {
			c.Abort()
			return response.Fail(c, 201, "Access Token不能为空")
		}

		token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.Common.TokenSecret), nil
		})

		if err != nil {
			c.Abort()
			return response.Fail(c, 201, "Access Token不合法")
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			user := &model.Users{}
			err = db.Model(user).Where("id = ?", claims.Issuer).Get()
			if err != nil {
				c.Abort()
				return response.Fail(c, 201, "Access Token错误，用户不存在")
			}

			c.Set("userInfo", user)
		} else {
			logger.Errorf("Access Token不合法 %s", err)
			c.Abort()
			return response.Fail(c, 201, "Access Token不合法")
		}

		c.Next()
		return nil
	}
}
