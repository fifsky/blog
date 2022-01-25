package handler

import (
	"net/http"
	"os"
	"testing"

	"app/config"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

func TestDingTalk_DingMsg(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments")...)
		conf := &config.Config{}
		conf.Common.DingAppSecret = "abc"
		handler := NewDingTalk(db, conf)

		t.Run("empty", func(t *testing.T) {
			req := test.NewRequest("/api/dingmsg", handler.DingMsg)
			req.Header.Set("timestamp", "2020-06-09 12:12:11")
			req.Header.Set("sign", "vLMGBUARoPek0ks3vv0eAi39mWkzS7afb+5njTlc7BY=")
			resp, err := req.JSON(``)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.Code)
			require.Equal(t, `{"msgtype":"text","text":{"content":"我还在开发哦……"}}`, resp.GetBodyString())
		})

		t.Run("sign error", func(t *testing.T) {
			req := test.NewRequest("/api/dingmsg", handler.DingMsg)
			req.Header.Set("timestamp", "2020-06-09 12:12:11")
			req.Header.Set("sign", "errrrrrrr")
			resp, err := req.JSON(`test`)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.Code)
			require.Equal(t, `{"msgtype":"text","text":{"content":"签名错误"}}`, resp.GetBodyString())
		})

		t.Run("mood", func(t *testing.T) {
			req := test.NewRequest("/api/dingmsg", handler.DingMsg)
			req.Header.Set("timestamp", "2020-06-09 12:12:11")
			req.Header.Set("sign", "vLMGBUARoPek0ks3vv0eAi39mWkzS7afb+5njTlc7BY=")
			resp, err := req.JSON(gee.H{"text": gee.H{"content": "#心情#开开心心每一天"}})
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.Code)
			require.Equal(t, `{"msgtype":"text","text":{"content":"心情发表成功"}}`, resp.GetBodyString())
		})

		t.Run("chatbot", func(t *testing.T) {
			if testing.Short() {
				t.Skip("skip")
			}

			conf.TencentCloud.SecretId = os.Getenv("TencentCloudSecretId")
			conf.TencentCloud.SecretKey = os.Getenv("TencentCloudSecretKey")

			req := test.NewRequest("/api/dingmsg", handler.DingMsg)
			req.Header.Set("timestamp", "2020-06-09 12:12:11")
			req.Header.Set("sign", "vLMGBUARoPek0ks3vv0eAi39mWkzS7afb+5njTlc7BY=")
			resp, err := req.JSON(gee.H{"text": gee.H{"content": "今天不开心"}})
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.Code)
			require.NotContains(t, `Ohoooo……`, resp.GetBodyString())
		})
	})
}
