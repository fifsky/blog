package service

import (
	"app/config"
	"app/pkg/httputil"
	apiv1 "app/proto/gen/api/v1"
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"
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
	if ok := verifyWeixinSignature(token, timestamp, nonce, expected); !ok {
		t.Fatalf("expect signature valid, got invalid")
	}

	// 错误签名应校验失败
	if ok := verifyWeixinSignature(token, timestamp, nonce, "wrong"); ok {
		t.Fatalf("expect signature invalid, got valid")
	}
}

func newTestHttpClient(isMock bool, suites []httputil.MockSuite) *http.Client {
	if isMock {
		return httputil.NewClient(httputil.WithMiddleware(httputil.Debug(), httputil.Mock(suites)))
	}
	return httputil.NewClient(httputil.WithMiddleware(httputil.Debug()))
}

func TestWeixin_getAccessToken(t *testing.T) {
	conf := &config.Config{}
	conf.Weixin.Appid = os.Getenv("WEIXIN_APPID")
	conf.Weixin.AppSecret = os.Getenv("WEIXIN_APPSECRET")
	conf.Weixin.Token = os.Getenv("WEIXIN_TOKEN")
	s := NewWeixin(nil, conf, newTestHttpClient(true, []httputil.MockSuite{
		{
			URI:          "/cgi-bin/stable_token",
			ResponseBody: `{"access_token":"mock_access_token","expires_in":7200}`,
		},
	}))
	resp, err := s.getAccessToken()
	if err != nil {
		t.Fatalf("getAccessToken failed: %v", err)
	}
	fmt.Println(resp.AccessToken)
	if resp.AccessToken == "" {
		t.Fatalf("expect access_token not empty, got empty")
	}
	if resp.ExpiresIn <= 0 {
		t.Fatalf("expect expires_in > 0, got %d", resp.ExpiresIn)
	}
}

func TestWeixin_Message(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	conf := &config.Config{}
	conf.Weixin.Appid = os.Getenv("WEIXIN_APPID")
	conf.Weixin.AppSecret = os.Getenv("WEIXIN_APPSECRET")
	conf.Weixin.Token = os.Getenv("WEIXIN_TOKEN")
	s := NewWeixin(nil, conf, httputil.NewClient(httputil.WithMiddleware(httputil.Debug())))
	req := &apiv1.MessageRequest{
		Touser:     "ofPKKs4q2T_Padf5oQU94JHESzkY",
		TemplateId: "mxrP6OFuLOy4UW1IKnsis0aR4FhxCXk77yS-LHdHvIQ",
		Data: map[string]*apiv1.MessageRequest_Value{
			"title":   {Value: "你好"},
			"content": {Value: "123456"},
		},
	}

	resp, err := s.Message(context.Background(), req)
	if err != nil {
		t.Fatalf("Message failed: %v", err)
	}
	fmt.Println(resp)
}
