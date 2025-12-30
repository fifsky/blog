package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"app/config"
	"app/pkg/aesutil"
	"app/store"
	"app/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/goapt/dbunit"
)

func getRemindTestToken(id int, conf *config.Config) string {
	token, _ := aesutil.AesEncode(conf.Common.TokenSecret, strconv.Itoa(id))
	return token
}

func TestNewRemindAuth(t *testing.T) {
	conf := &config.Config{}
	conf.Common.TokenSecret = "abcdabcdabcdabcd"

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("reminds"))

		tests := []struct {
			Token        string
			ResponseBody string
		}{
			{
				getRemindTestToken(8, conf),
				`{"code":10000,"msg":"success"}`,
			},
			{
				"",
				`{"code":201,"msg":"非法访问"}`,
			},
			{
				"789789",
				`{"code":202,"msg":"Token错误"}`,
			},
			{
				getRemindTestToken(888, conf),
				`{"code":203,"msg":"数据不存在"}`,
			},
		}

		// 构造下游处理器：仅在通过鉴权时返回成功JSON
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code":10000,"msg":"success"}`))
		})

		for _, tt := range tests {
			// 创建被测的中间件
			m := NewRemindAuth(store.New(db), conf)

			// 包装中间件
			handler := m(next)

			// 构造请求，设置查询参数token
			req := httptest.NewRequest(http.MethodGet, "/remind?token="+tt.Token, nil)
			// 记录响应
			rr := httptest.NewRecorder()

			// 触发请求
			handler.ServeHTTP(rr, req)

			assert.JSONEq(t, tt.ResponseBody, rr.Body.String())
		}
	})
}
