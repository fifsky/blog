package middleware

import (
	"fmt"
	"testing"

	"app/config"
	"app/pkg/aesutil"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/golib/hashing"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
)

func getUserTestToken(id int, conf *config.Config) string {
	src := fmt.Sprintf("%d:%s", id, hashing.Md5(fmt.Sprintf("%d%s", id, conf.Common.TokenSecret)))

	if id == 888 {
		src = src + "err"
	}

	token, _ := aesutil.AesEncode(conf.Common.TokenSecret, src)
	return token
}

func TestNewAuthLogin(t *testing.T) {
	var testHandler gee.HandlerFunc = func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"msg":  "success",
		})
	}

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
				`{"code":201,"msg":"Access Token错误"}`,
			},
			{
				getUserTestToken(888, conf),
				`{"code":201,"msg":"Access Token不合法"}`,
			},
			{
				getUserTestToken(999, conf),
				`{"code":201,"msg":"Access Token错误，用户不存在"}`,
			},
		}

		for _, tt := range tests {
			req := test.NewRequest("/dummy/impl", gee.HandlerFunc(NewAuthLogin(db, conf)), testHandler)
			req.Header.Add("Access-Token", tt.Token)
			resp, err := req.Get()
			assert.NoError(t, err)
			assert.Equal(t, tt.ResponseBody, resp.GetBodyString())
		}
	})

}
