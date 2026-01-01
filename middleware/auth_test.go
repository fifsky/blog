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
			{
				getUserTestToken(1, conf),
				`{"code":10000,"msg":"success"}`,
			},
			{
				"",
				`{"code":201,"msg":"Access Token不能为空"}`,
			},
			{
				"789789",
				`{"code":201,"msg":"Access Token不合法"}`,
			},
			{
				getUserTestToken(999, conf),
				`{"code":202,"msg":"Access Token错误，用户不存在"}`,
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
