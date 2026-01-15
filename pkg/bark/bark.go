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
	Title    string `json:"title"`
	Body     string `json:"body,omitempty"`
	Badge    int    `json:"badge,omitempty"`
	Url      string `json:"url,omitempty"`
	Markdown string `json:"markdown,omitempty"`
	Group    string `json:"group,omitempty"`
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
