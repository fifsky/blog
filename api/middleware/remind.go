package middleware

import (
	"net/http"
	"strconv"

	"app/handler"
	"app/response"
	"app/store"

	"app/config"
	"app/pkg/aesutil"

	"github.com/samber/lo"
)

type RemindAuth = func(next http.Handler) http.Handler

func NewRemindAuth(s *store.Store, conf *config.Config) RemindAuth {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")

			if token == "" {
				response.Fail(w, 201, "非法访问")
				return
			}

			id, err := aesutil.AesDecode(conf.Common.TokenSecret, token)
			if err != nil {
				response.Fail(w, 202, "Token错误")
				return
			}

			remind, err := s.GetRemind(r.Context(), lo.Must(strconv.Atoi(id)))
			if err != nil {
				response.Fail(w, 203, "数据不存在")
				return
			}

			r = r.WithContext(handler.SetRemind(r.Context(), remind))
			next.ServeHTTP(w, r)
		})
	}
}
