package handler

import (
	"net/http"
	"testing"

	"app/config"
	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

func TestComment_Avatar(t *testing.T) {
	handler := NewCommon()
	req := test.NewRequest("/api/avatar", handler.Avatar)
	resp, err := req.JSON(``)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)
}

func TestCommon_DingMsg(t *testing.T) {
	config.App.Common.DingAppSecret = "abc"
	t.Run("success", func(t *testing.T) {
		handler := NewCommon()
		req := test.NewRequest("/api/dingmsg", handler.DingMsg)
		req.Header.Set("timestamp", "2020-06-09 12:12:11")
		req.Header.Set("sign", "vLMGBUARoPek0ks3vv0eAi39mWkzS7afb+5njTlc7BY=")
		resp, err := req.JSON(`test`)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, `{"msgtype":"text","text":{"content":"我还在开发哦……"}}`, resp.GetBodyString())
	})

	t.Run("sign error", func(t *testing.T) {
		handler := NewCommon()
		req := test.NewRequest("/api/dingmsg", handler.DingMsg)
		req.Header.Set("timestamp", "2020-06-09 12:12:11")
		req.Header.Set("sign", "errrrrrrr")
		resp, err := req.JSON(`test`)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, `{"code":202,"msg":"签名错误"}`, resp.GetBodyString())
	})
}
