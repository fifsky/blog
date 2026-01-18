package bark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	url        string
	token      string
}

func New(client *http.Client, url, token string) *Client {
	return &Client{
		httpClient: client,
		url:        url,
		token:      token,
	}
}

type Message struct {
	// 推送标题
	Title string `json:"title"`
	// 推送内容
	Body string `json:"body,omitempty"`
	// 推送角标，可以是任意数字
	Badge int `json:"badge,omitempty"`
	//推送中断级别。
	// critical: 重要警告, 在静音模式下也会响铃
	// active：默认值，系统会立即亮屏显示通知
	// timeSensitive：时效性通知，可在专注状态下显示通知。
	// passive：仅将通知添加到通知列表，不会亮屏提醒。
	Level string `json:"level,omitempty"`
	// 点击通知后跳转的 URL
	Url string `json:"url,omitempty"`
	// 推送内容是否支持 Markdown 格式
	Markdown string `json:"markdown,omitempty"`
	// 推送分组，用于对通知进行分类管理
	Group string `json:"group,omitempty"`
}

func (c *Client) Send(msg Message) error {
	reqBody, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal bark message error: %w", err)
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if c.token != "" {
		req.Header.Set("Authorization", "Basic "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("bark api returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %w", err)
	}
	fmt.Println(string(body))

	return nil
}
