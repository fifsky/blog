package middleware

import (
	"net/http"
	"strconv"

	"app/config"
	"app/pkg/errors"
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
				response.Fail(w, errors.ErrUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(conf.Common.TokenSecret), nil
			})

			if err != nil {
				response.Fail(w, errors.ErrUnauthorized.WithCause(err))
				return
			}

			if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
				user, err := s.GetUser(r.Context(), lo.Must(strconv.Atoi(claims.Issuer)))
				if err != nil {
					response.Fail(w, errors.ErrUnauthorized.WithMetadata(map[string]string{"user_id": claims.Issuer}))
					return
				}

				r = r.WithContext(admin.SetLoginUser(r.Context(), user))
			} else {
				response.Fail(w, errors.ErrUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
