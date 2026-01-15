package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"app/pkg/errors"
	"app/response"

	"github.com/goapt/logger"
)

func NewRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				stack := string(buf[:n])

				// 记录错误日志
				logger.Default().Error(fmt.Sprintf("%v", err), "stack", stack)

				// 返回系统错误响应
				response.Fail(w, errors.ErrSystem.WithCause(fmt.Errorf("%v", err)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
