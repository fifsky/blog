package service

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"app/config"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
)

var _ apiv1.WeixinServiceServer = (*Weixin)(nil)

// Weixin 使用微信测试账号发送消息 https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index
type Weixin struct {
	apiv1.UnimplementedWeixinServiceServer
	store      *store.Store
	conf       *config.Config
	httpClient *http.Client
}

// NewWeixin 构造函数
func NewWeixin(store *store.Store, conf *config.Config, httpClient *http.Client) *Weixin {
	return &Weixin{store: store, conf: conf, httpClient: httpClient}
}

// verifyWeixinSignature 按微信官方文档校验签名
// 校验流程：
// 1) 将 token、timestamp、nonce 三个参数进行字典序排序
// 2) 将三个参数字符串拼接成一个字符串进行 sha1 加密
// 3) 与请求中的 signature 比对，一致则校验通过
func verifyWeixinSignature(token, timestamp, nonce, signature string) bool {
	// 将三个参数放入切片并按字典序排序
	parts := []string{token, timestamp, nonce}
	sort.Strings(parts)

	// 拼接为单个字符串
	raw := strings.Join(parts, "")

	// 计算 sha1 摘要，并转为十六进制小写字符串
	sum := sha1.Sum([]byte(raw))
	sign := fmt.Sprintf("%x", sum)

	// 比较计算结果与微信传入的签名
	return sign == signature
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
}

type ErrorResponse struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// 使用 sync.Map 作为全局缓存：key=appid，value=cacheEntry
type cacheEntry struct {
	resp      *AccessTokenResponse
	expiresAt time.Time
}

var tokenCache sync.Map

func (s *Weixin) getAccessToken() (*AccessTokenResponse, error) {
	// 读取配置的 appid、secret
	req := &AccessTokenRequest{
		GrantType: "client_credential",
		Appid:     s.conf.Weixin.Appid,
		Secret:    s.conf.Weixin.AppSecret,
	}

	// 命中未过期的缓存则直接返回
	// 使用 appid 作为缓存键（一个应用仅依赖 appid）
	key := req.Appid
	if v, ok := tokenCache.Load(key); ok {
		entry := v.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.resp, nil
		}
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://api.weixin.qq.com/cgi-bin/stable_token", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 先解析到错误响应
	errResp := &ErrorResponse{}
	err = json.Unmarshal(body, errResp)
	if err != nil {
		return nil, err
	}
	if errResp.Errcode != 0 {
		return nil, fmt.Errorf("error response: %d, %s", errResp.Errcode, errResp.Errmsg)
	}

	// 再解析到具体的结构体
	respBody := &AccessTokenResponse{}
	err = json.Unmarshal(body, respBody)
	if err != nil {
		return nil, err
	}
	// 计算过期时间，预留120秒安全缓冲避免临界点失效
	ttl := time.Duration(respBody.ExpiresIn) * time.Second
	if ttl <= 0 {
		// 若返回异常，设置一个默认短TTL确保后续可刷新
		ttl = 5 * time.Minute
	}
	expiresAt := time.Now().Add(ttl - 120*time.Second)

	// 写入缓存
	tokenCache.Store(key, cacheEntry{resp: respBody, expiresAt: expiresAt})

	return respBody, nil
}

func (s *Weixin) Notify(w http.ResponseWriter, r *http.Request) {
	// 统一读取查询参数
	q := r.URL.Query()
	signature := q.Get("signature")
	timestamp := q.Get("timestamp")
	nonce := q.Get("nonce")
	echostr := q.Get("echostr")

	// 从环境或配置获取开发者在微信平台配置的 Token
	token := s.conf.Weixin.Token

	// 先进行签名校验，失败则直接返回 403
	if !verifyWeixinSignature(token, timestamp, nonce, signature) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("invalid signature"))
		return
	}

	// GET 请求：用于接入时的 URL 有效性验证，需原样返回 echostr
	if r.Method == http.MethodGet {
		_, _ = w.Write([]byte(echostr))
		return
	}

	// POST 请求：微信的消息与事件推送，签名已通过
	// 此处预留解析 XML 及业务处理逻辑
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Message 发送模板消息
func (s *Weixin) Message(ctx context.Context, req *apiv1.MessageRequest) (*apiv1.WeixinResponse, error) {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken.AccessToken)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析响应体
	respBody := &ErrorResponse{}
	err = json.Unmarshal(body, respBody)
	if err != nil {
		return nil, err
	}
	if respBody.Errcode != 0 {
		return nil, fmt.Errorf("error response: %d, %s", respBody.Errcode, respBody.Errmsg)
	}

	return &apiv1.WeixinResponse{}, nil
}
