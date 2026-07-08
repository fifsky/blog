// Package doubaoasr 封装豆包语音识别极速版 HTTP 接口。
package doubaoasr

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultEndpoint   = "https://openspeech.bytedance.com/api/v3/auc/bigmodel/recognize/flash"
	defaultResourceID = "volc.bigasr.auc_turbo"
	defaultModelName  = "bigmodel"
	successCodePrefix = "200"
)

// Client 是豆包语音识别 HTTP 客户端。
type Client struct {
	// APIKey 是新版控制台提供的 X-Api-Key。
	APIKey string
	// Endpoint 是极速版识别接口地址，留空使用默认地址。
	Endpoint string
	// ResourceID 是豆包语音资源 ID，留空使用默认极速版资源。
	ResourceID string
	// UID 是透传给豆包的用户标识，留空使用应用默认值。
	UID string
	// HTTPClient 是可注入的 HTTP 客户端，测试时用于替换网络调用。
	HTTPClient *http.Client
}

// Transcribe 将 base64 音频发送到豆包语音识别接口并返回文字。
func (c Client) Transcribe(ctx context.Context, audioBase64 string) (string, error) {
	apiKey := strings.TrimSpace(c.APIKey)
	if apiKey == "" {
		return "", errors.New("doubao asr api key is empty")
	}
	audioBase64 = strings.TrimSpace(audioBase64)
	if audioBase64 == "" {
		return "", errors.New("audio base64 is empty")
	}

	payload := recognizeRequest{
		User: recognizeUser{UID: c.uid()},
		Audio: recognizeAudio{
			Data: audioBase64,
		},
		Request: recognizeOptions{
			ModelName: defaultModelName,
		},
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(payload); err != nil {
		return "", fmt.Errorf("encode doubao asr request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), &body)
	if err != nil {
		return "", fmt.Errorf("create doubao asr request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("X-Api-Resource-Id", c.resourceID())
	req.Header.Set("X-Api-Request-Id", newRequestID())
	req.Header.Set("X-Api-Sequence", "-1")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("call doubao asr: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read doubao asr response: %w", err)
	}

	headers := responseHeaders{
		APIStatusCode: strings.TrimSpace(resp.Header.Get("X-Api-Status-Code")),
		APIMessage:    strings.TrimSpace(resp.Header.Get("X-Api-Message")),
		LogID:         strings.TrimSpace(resp.Header.Get("X-Tt-Logid")),
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", responseError(resp.Status, headers, data)
	}
	if headers.APIStatusCode != "" && !strings.HasPrefix(headers.APIStatusCode, successCodePrefix) {
		return "", responseError(resp.Status, headers, data)
	}

	var result recognizeResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("decode doubao asr response: %w", err)
	}

	text := strings.TrimSpace(result.Result.Text)
	if text == "" {
		return "", emptyTextError(headers, data)
	}
	return text, nil
}

func (c Client) endpoint() string {
	if endpoint := strings.TrimSpace(c.Endpoint); endpoint != "" {
		return endpoint
	}
	return defaultEndpoint
}

func (c Client) resourceID() string {
	if resourceID := strings.TrimSpace(c.ResourceID); resourceID != "" {
		return resourceID
	}
	return defaultResourceID
}

func (c Client) uid() string {
	if uid := strings.TrimSpace(c.UID); uid != "" {
		return uid
	}
	return "fifsky-blog-ios"
}

func (c Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 60 * time.Second}
}

func newRequestID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return fmt.Sprintf("blog-%d", time.Now().UnixNano())
	}
	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80
	parts := []string{
		hex.EncodeToString(raw[0:4]),
		hex.EncodeToString(raw[4:6]),
		hex.EncodeToString(raw[6:8]),
		hex.EncodeToString(raw[8:10]),
		hex.EncodeToString(raw[10:16]),
	}
	return strings.Join(parts, "-")
}

func responseError(httpStatus string, headers responseHeaders, body []byte) error {
	parts := []string{"doubao asr failed"}
	if strings.TrimSpace(httpStatus) != "" {
		parts = append(parts, "http_status="+strings.TrimSpace(httpStatus))
	}
	if headers.APIStatusCode != "" {
		parts = append(parts, "api_status_code="+headers.APIStatusCode)
	}
	if headers.APIMessage != "" {
		parts = append(parts, "api_message="+headers.APIMessage)
	}
	if headers.LogID != "" {
		parts = append(parts, "logid="+headers.LogID)
	}
	if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
		parts = append(parts, "body="+truncate(trimmed, 240))
	}
	return errors.New(strings.Join(parts, "; "))
}

func emptyTextError(headers responseHeaders, body []byte) error {
	parts := []string{"doubao asr returned empty text"}
	if headers.APIStatusCode != "" {
		parts = append(parts, "api_status_code="+headers.APIStatusCode)
	}
	if headers.APIMessage != "" {
		parts = append(parts, "api_message="+headers.APIMessage)
	}
	if headers.LogID != "" {
		parts = append(parts, "logid="+headers.LogID)
	}
	if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
		parts = append(parts, "raw="+truncate(trimmed, 240))
	}
	return errors.New(strings.Join(parts, "; "))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

type recognizeRequest struct {
	User    recognizeUser    `json:"user"`
	Audio   recognizeAudio   `json:"audio"`
	Request recognizeOptions `json:"request"`
}

type recognizeUser struct {
	UID string `json:"uid"`
}

type recognizeAudio struct {
	Data string `json:"data"`
}

type recognizeOptions struct {
	ModelName string `json:"model_name"`
}

type recognizeResponse struct {
	Result recognizeResult `json:"result"`
}

type recognizeResult struct {
	Text string `json:"text"`
}

type responseHeaders struct {
	APIStatusCode string
	APIMessage    string
	LogID         string
}
