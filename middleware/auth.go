package middleware

import (
	"net/http"
	"strconv"

	"app/config"
	"app/response"
	"app/service/admin"
	"app/store"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

type AuthLogin = func(next http.Handler) http.Handler

func NewAuthLogin(s *store.Store, conf *config.Config) AuthLogin {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Access-Token")

			if accessToken == "" {
				response.Fail(w, 201, "登录过期，请重新登录")
				return
			}

			token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(conf.Common.TokenSecret), nil
			})

			if err != nil {
				response.Fail(w, 201, "登录过期，请重新登录")
				return
			}

			if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
				user, err := s.GetUser(r.Context(), lo.Must(strconv.Atoi(claims.Issuer)))
				if err != nil {
					response.Fail(w, 202, "Access Token错误，用户不存在")
					return
				}

				r = r.WithContext(admin.SetLoginUser(r.Context(), user))
			} else {
				response.Fail(w, 203, "登录过期，请重新登录")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
