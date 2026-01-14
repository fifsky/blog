package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"app/response"
	"github.com/goapt/logger"
)

func NewRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Println(string(buf))
				logger.Default().Error(fmt.Sprintf("%v", err), "stack", string(buf))
				response.Fail(w, http.StatusInternalServerError, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
