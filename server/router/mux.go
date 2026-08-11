package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
)

type ServeMux struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux: http.NewServeMux(),
	}
}

// Use 原地追加中间件到当前 mux，影响后续注册的所有 handler
func (s *ServeMux) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Group 创建新的 ServeMux，共享底层 http.ServeMux，中间件栈继承父级并叠加新增中间件，
// 仅影响分组内注册的 handler，不修改原 mux
func (s *ServeMux) Group(middlewares ...Middleware) *ServeMux {
	return &ServeMux{
		ServeMux:    s.ServeMux,
		middlewares: append(s.middlewares, middlewares...),
	}
}

func (s *ServeMux) Handle(pattern string, handler http.Handler) {
	slog.Info(fmt.Sprintf("[http] %-25s  --> %s", pattern, nameOfFunction(handler)))
	s.ServeMux.Handle(pattern, chain(handler, s.middlewares))
}

func (s *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.Handle(pattern, http.HandlerFunc(handler))
}

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
