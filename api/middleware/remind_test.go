package middleware

import (
	"strconv"
	"testing"

	"app/config"
	"app/pkg/aesutil"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
)

func getRemindTestToken(id int) string {
	token, _ := aesutil.AesEncode(config.App.Common.TokenSecret, strconv.Itoa(id))
	return token
}

func TestNewRemindAuth(t *testing.T) {
	var testHandler gee.HandlerFunc = func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"msg":  "success",
		})
	}

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("reminds"))

		tests := []struct {
			Token        string
			ResponseBody string
		}{
			{
				getRemindTestToken(8),
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
				getRemindTestToken(888),
				`{"code":203,"msg":"数据不存在"}`,
			},
		}

		for _, tt := range tests {
			req := test.NewRequest("/dummy/impl?token="+tt.Token, gee.HandlerFunc(NewRemindAuth(db)), testHandler)
			resp, err := req.Get()
			assert.NoError(t, err)
			assert.Equal(t, tt.ResponseBody, resp.GetBodyString())
		}
	})
}
