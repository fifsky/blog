package openapi

import (
	"app/config"
	apiv1 "app/proto/gen/api/v1"
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/goapt/httpx"
	"github.com/stretchr/testify/assert"
)

// TestVerifyWeixinSignature 验证签名校验流程是否符合微信官方要求
func TestVerifyWeixinSignature(t *testing.T) {
	// 模拟开发者在微信平台配置的 Token
	token := "mytoken"
	// 微信请求携带的时间戳与随机数（字符串形式）
	timestamp := "1736680150"
	nonce := "abcdef"

	// 按微信规则计算期望的签名：对 token、timestamp、nonce 字典序排序并拼接后 sha1
	parts := []string{token, timestamp, nonce}
	sort.Strings(parts)
	raw := strings.Join(parts, "")
	expected := fmt.Sprintf("%x", sha1.Sum([]byte(raw)))

	// 正确签名应校验通过
	assert.True(t, verifyWeixinSignature(token, timestamp, nonce, expected), "expect signature valid, got invalid")

	// 错误签名应校验失败
	assert.False(t, verifyWeixinSignature(token, timestamp, nonce, "wrong"), "expect signature invalid, got valid")
}

func newTestHttpClient(isMock bool, suites []httpx.MockSuite) *http.Client {
	if isMock {
		return httpx.NewClient(httpx.WithMiddleware(httpx.Debug(), httpx.Mock(suites)))
	}
	return httpx.NewClient(httpx.WithMiddleware(httpx.Debug()))
}

func TestWeixin_getAccessToken(t *testing.T) {
	conf := &config.Config{}
	conf.Weixin.Appid = os.Getenv("WEIXIN_APPID")
	conf.Weixin.AppSecret = os.Getenv("WEIXIN_APPSECRET")
	conf.Weixin.Token = os.Getenv("WEIXIN_TOKEN")
	s := NewWeixin(nil, conf, newTestHttpClient(true, []httpx.MockSuite{
		{
			URI:          "/cgi-bin/stable_token",
			ResponseBody: `{"access_token":"mock_access_token","expires_in":7200}`,
		},
	}))
	resp, err := s.getAccessToken()
	assert.NoError(t, err)
	t.Log(resp.AccessToken)
	assert.NotEmpty(t, resp.AccessToken, "expect access_token not empty")
	assert.Positive(t, resp.ExpiresIn, "expect expires_in > 0")
}

func TestWeixin_Message(t *testing.T) {
	conf := &config.Config{}
	conf.Weixin.Appid = os.Getenv("WEIXIN_APPID")
	conf.Weixin.AppSecret = os.Getenv("WEIXIN_APPSECRET")
	conf.Weixin.Token = os.Getenv("WEIXIN_TOKEN")
	s := NewWeixin(nil, conf, newTestHttpClient(true, []httpx.MockSuite{
		{
			URI:          "/cgi-bin/stable_token",
			ResponseBody: `{"access_token":"mock_access_token","expires_in":7200}`,
		},
		{
			URI:          "/cgi-bin/message/template/send",
			ResponseBody: `{"errcode":0,"errmsg":"ok"}`,
		},
	}))
	req := &apiv1.MessageRequest{
		Touser:     "ofPKKs4q2T_Padf5oQU94JHESzkY",
		TemplateId: "mxrP6OFuLOy4UW1IKnsis0aR4FhxCXk77yS-LHdHvIQ",
		Data: map[string]*apiv1.MessageRequest_Value{
			"title":   {Value: "你好"},
			"content": {Value: "123456"},
		},
	}

	resp, err := s.Message(context.Background(), req)
	assert.NoError(t, err)
	t.Log(resp)
}
