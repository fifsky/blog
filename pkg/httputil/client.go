package httputil

import (
	"net/http"
	"time"

	"app/pkg/httputil/middleware"
)

type Option func(c *Client)

type Client struct {
	timeout     time.Duration
	middlewares []middleware.Middleware
	transport   http.RoundTripper
}

func WithTimeout(t time.Duration) Option {
	return func(c *Client) {
		c.timeout = t
	}
}

func WithMiddleware(middleware ...middleware.Middleware) Option {
	return func(c *Client) {
		c.middlewares = append(c.middlewares, middleware...)
	}
}

func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.transport = transport
	}
}

func NewClient(options ...Option) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	// 这里的关键是：IdleConnTimeout 必须小于云厂商 LB 的超时时间 (通常是60s)
	// 设置为 30s 是比较稳妥的选择
	t.IdleConnTimeout = 30 * time.Second
	// 根据你的并发量调整，默认是2，太小会导致连接频繁创建销毁
	t.MaxIdleConnsPerHost = 10

	c := &Client{
		timeout:     5 * time.Second,
		middlewares: []middleware.Middleware{},
		transport:   t,
	}
	for _, option := range options {
		option(c)
	}

	return &http.Client{
		Transport: middleware.Chain(c.transport, c.middlewares...),
		Timeout:   c.timeout,
	}
}
