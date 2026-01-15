package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/config"
	"app/store"
	"app/testutil"

	"github.com/stretchr/testify/assert"

	"github.com/goapt/dbunit"
	"github.com/golang-jwt/jwt/v5"
)

func getUserTestToken(id int, conf *config.Config) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Issuer:    fmt.Sprint(id),
	})

	tokenString, _ := token.SignedString([]byte(conf.Common.TokenSecret))
	// fmt.Println("===>", tokenString, err)
	return tokenString
}

func TestNewAuthLogin(t *testing.T) {
	conf := &config.Config{}
	conf.Common.TokenSecret = "abcdabcdabcdabcd"

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))

		tests := []struct {
			Token        string
			ResponseBody string
		}{
			// 1. 正常登录，Token 有效且用户存在
			{
				getUserTestToken(1, conf),
				`{"code":10000,"msg":"success"}`,
			},
			// 2. 未提供 Token
			{
				"",
				`{"code":"UNAUTHORIZED","message":"登录过期，请重新登录"}`,
			},
			// 3. Token 格式错误
			{
				"789789",
				`{"code":"UNAUTHORIZED","message":"登录过期，请重新登录","details":{"cause":"token is malformed: token contains an invalid number of segments"}}`,
			},
			// 4. Token 有效但用户不存在
			{
				getUserTestToken(999, conf),
				`{"code":"UNAUTHORIZED","message":"登录过期，请重新登录","details":{"user_id":"999"}}`,
			},
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code":10000,"msg":"success"}`))
		})

		for _, tt := range tests {
			m := NewAuthLogin(store.New(db), conf)

			handler := m(next)

			req := httptest.NewRequest(http.MethodGet, "/dummy/impl", nil)
			req.Header.Set("Access-Token", tt.Token)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.JSONEq(t, tt.ResponseBody, rr.Body.String())
		}
	})

}
