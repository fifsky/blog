package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	defaultStopTimeout       = 10 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 30 * time.Second
	defaultWriteTimeout      = 30 * time.Second
	defaultIdleTimeout       = 2 * time.Minute
)

// Option 表示 HTTP 服务配置项。
type Option func(*Server)

// Handler 设置 HTTP 请求处理器。
func Handler(handler http.Handler) Option {
	return func(s *Server) {
		s.Handler = handler
	}
}

// Address 设置 HTTP 服务监听地址。
func Address(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

// StopTimeout 设置 HTTP 服务优雅关闭的最大等待时间。
func StopTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.stopTimeout = t
	}
}

// ReadHeaderTimeout 设置读取请求头的最大等待时间。
func ReadHeaderTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.ReadHeaderTimeout = t
	}
}

// IdleTimeout 设置空闲连接的最大等待时间。
func IdleTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.IdleTimeout = t
	}
}

// ReadTimeout 设置读取完整请求的最大等待时间。
func ReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.ReadTimeout = t
	}
}

// WriteTimeout 设置写入响应的最大等待时间。
func WriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.WriteTimeout = t
	}
}

// Server 封装 HTTP 服务及其优雅关闭配置。
type Server struct {
	http.Server
	stopTimeout time.Duration
}

// New 创建 HTTP 服务并应用默认超时配置。
func New(opts ...Option) *Server {
	srv := &Server{
		Server: http.Server{
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			ReadTimeout:       defaultReadTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
		},
		stopTimeout: defaultStopTimeout,
	}

	for _, o := range opts {
		o(srv)
	}

	return srv
}

// Start 启动 HTTP 服务，并在运行上下文取消后执行优雅关闭。
func (s *Server) Start(ctx context.Context) error {
	log.Println("[HTTP] Server listen:" + s.Addr)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.serve()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		shutdownErr := s.shutdown()
		return errors.Join(shutdownErr, <-serveErr)
	}
}

// Stop 使用调用方提供的上下文关闭 HTTP 服务。
func (s *Server) Stop(ctx context.Context) error {
	log.Println("[HTTP] server stopping")
	return s.Shutdown(ctx)
}

func (s *Server) serve() error {
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server listen: %w", err)
	}
	return nil
}

func (s *Server) shutdown() error {
	log.Printf("[HTTP] Shutdown timeout: %s\n", s.stopTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", errors.Join(err, s.Close()))
	}
	return nil
}
