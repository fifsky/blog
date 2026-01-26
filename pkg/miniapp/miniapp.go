package miniapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Appid      string
	AppSecret  string
	HTTPClient *http.Client
}

type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func NewClient(appid, secret string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	return &Client{
		Appid:      appid,
		AppSecret:  secret,
		HTTPClient: httpClient,
	}
}

func (c *Client) Code2Session(ctx context.Context, code string) (*Code2SessionResponse, error) {
	values := url.Values{}
	values.Set("appid", c.Appid)
	values.Set("secret", c.AppSecret)
	values.Set("js_code", code)
	values.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weixin.qq.com/sns/jscode2session?"+values.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weixin error: %w", err)
	}
	defer resp.Body.Close()

	var ret Code2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("weixin error: %d %s", ret.ErrCode, ret.ErrMsg)
	}

	if ret.OpenID == "" {
		return nil, fmt.Errorf("weixin response openid empty")
	}

	return &ret, nil
}
