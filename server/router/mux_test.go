package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 以下中间件使用具名函数：闭包的代码指针相同，无法通过 reflect 断言身份，
// 而 Group 的核心风险正是"中间件被替换"，需要能区分出到底是哪一个
func mwBase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Base", "base")
		next.ServeHTTP(w, r)
	})
}

func mwFirst(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Group", "first")
		next.ServeHTTP(w, r)
	})
}

func mwSecond(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Group", "second")
		next.ServeHTTP(w, r)
	})
}

// middlewarePointer 返回中间件函数指针，用于断言中间件身份
func middlewarePointer(mw Middleware) uintptr {
	return reflect.ValueOf(mw).Pointer()
}

// newMuxWithSpareCap 构造父级中间件栈存在富余容量的 mux，
// 复现 append 原地写入底层数组的场景
func newMuxWithSpareCap() *ServeMux {
	mux := &ServeMux{ServeMux: http.NewServeMux()}
	mux.middlewares = append(make([]Middleware, 0, 8), mwBase)
	return mux
}

func TestServeMux_Group(t *testing.T) {
	t.Run("兄弟分组中间件互不覆盖", func(t *testing.T) {
		mux := newMuxWithSpareCap()

		a := mux.Group(mwFirst)
		b := mux.Group(mwSecond)

		assert.Len(t, a.middlewares, 2)
		assert.Len(t, b.middlewares, 2)
		assert.Equal(t, middlewarePointer(mwFirst), middlewarePointer(a.middlewares[1]))
		assert.Equal(t, middlewarePointer(mwSecond), middlewarePointer(b.middlewares[1]))
	})

	t.Run("分组继承父级中间件且不修改父级", func(t *testing.T) {
		mux := newMuxWithSpareCap()

		child := mux.Group(mwFirst)

		assert.Len(t, mux.middlewares, 1)
		assert.Equal(t, middlewarePointer(mwBase), middlewarePointer(mux.middlewares[0]))
		assert.Len(t, child.middlewares, 2)
		assert.Equal(t, middlewarePointer(mwBase), middlewarePointer(child.middlewares[0]))
	})

	t.Run("后创建的分组不影响先创建分组的handler链", func(t *testing.T) {
		mux := newMuxWithSpareCap()

		a := mux.Group(mwFirst)
		b := mux.Group(mwSecond)
		// 分组与注册分离：先建好两个分组，再分别注册路由
		a.HandleFunc("GET /first", func(w http.ResponseWriter, r *http.Request) {})
		b.HandleFunc("GET /second", func(w http.ResponseWriter, r *http.Request) {})

		firstRec := httptest.NewRecorder()
		a.ServeHTTP(firstRec, httptest.NewRequest(http.MethodGet, "/first", nil))
		assert.Equal(t, "first", firstRec.Header().Get("X-Group"))
		assert.Equal(t, "base", firstRec.Header().Get("X-Base"))

		secondRec := httptest.NewRecorder()
		b.ServeHTTP(secondRec, httptest.NewRequest(http.MethodGet, "/second", nil))
		assert.Equal(t, "second", secondRec.Header().Get("X-Group"))
		assert.Equal(t, "base", secondRec.Header().Get("X-Base"))
	})
}
