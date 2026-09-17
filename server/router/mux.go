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
// 仅影响分组内注册的 handler，不修改原 mux。
// 这里必须复制到新切片：若父级切片仍有富余容量，append 会原地写入底层数组，
// 导致多个兄弟分组共享同一底层数组而相互覆盖
func (s *ServeMux) Group(middlewares ...Middleware) *ServeMux {
	ms := make([]Middleware, 0, len(s.middlewares)+len(middlewares))
	ms = append(ms, s.middlewares...)
	ms = append(ms, middlewares...)
	return &ServeMux{
		ServeMux:    s.ServeMux,
		middlewares: ms,
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
