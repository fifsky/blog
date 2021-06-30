package middleware

import (
	"testing"

	"app/config"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Cors(t *testing.T) {
	var testHandler gee.HandlerFunc = func(c *gee.Context) gee.Response {
		return c.JSON(gee.H{
			"code": 10000,
			"msg":  "success",
		})
	}

	{
		req := test.NewRequest("/dummy/impl", gee.HandlerFunc(NewCors()), testHandler)
		req.Host = "fifsky.com"
		req.Header.Set("Origin", "http://fifsky.com")

		resp, err := req.JSON(``)
		assert.NoError(t, err)
		assert.Equal(t, `{"code":10000,"msg":"success"}`, resp.GetBodyString())
	}

	{
		config.App.Env = "local"
		req := test.NewRequest("/dummy/impl", gee.HandlerFunc(NewCors()), testHandler)
		req.Host = "test.com"
		req.Header.Set("Origin", "http://test.com")

		resp, err := req.JSON(``)
		assert.NoError(t, err)
		assert.Equal(t, `{"code":10000,"msg":"success"}`, resp.GetBodyString())
	}
}
