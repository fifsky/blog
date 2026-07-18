package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name              string
		opts              []Option
		stopTimeout       time.Duration
		readHeaderTimeout time.Duration
		readTimeout       time.Duration
		writeTimeout      time.Duration
		idleTimeout       time.Duration
	}{
		{
			name:              "默认配置",
			stopTimeout:       defaultStopTimeout,
			readHeaderTimeout: defaultReadHeaderTimeout,
			readTimeout:       defaultReadTimeout,
			writeTimeout:      defaultWriteTimeout,
			idleTimeout:       defaultIdleTimeout,
		},
		{
			name: "覆盖超时配置",
			opts: []Option{
				StopTimeout(time.Second),
				ReadHeaderTimeout(2 * time.Second),
				ReadTimeout(3 * time.Second),
				WriteTimeout(4 * time.Second),
				IdleTimeout(5 * time.Second),
			},
			stopTimeout:       time.Second,
			readHeaderTimeout: 2 * time.Second,
			readTimeout:       3 * time.Second,
			writeTimeout:      4 * time.Second,
			idleTimeout:       5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := New(tt.opts...)

			assert.Equal(t, tt.stopTimeout, srv.stopTimeout)
			assert.Equal(t, tt.readHeaderTimeout, srv.ReadHeaderTimeout)
			assert.Equal(t, tt.readTimeout, srv.ReadTimeout)
			assert.Equal(t, tt.writeTimeout, srv.WriteTimeout)
			assert.Equal(t, tt.idleTimeout, srv.IdleTimeout)
		})
	}
}

func TestServer_StartWithCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	srv := New(
		Address("127.0.0.1:0"),
		Handler(http.NotFoundHandler()),
	)
	require.NoError(t, srv.Start(ctx))
}

func TestServer_StartListenError(t *testing.T) {
	srv := New(Address("127.0.0.1:-1"))
	err := srv.Start(context.Background())
	require.ErrorContains(t, err, "http server listen")
}

func TestServer_StartGracefulShutdown(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{}, 1)
	defer close(releaseRequest)
	srv := New(
		Address(addr),
		Handler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			close(requestStarted)
			<-releaseRequest
		})),
		StopTimeout(time.Second),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startErr := make(chan error, 1)
	go func() {
		startErr <- srv.Start(ctx)
	}()

	waitForServer(t, addr)
	requestDone := requestServer(t, addr)
	waitForSignal(t, requestStarted, "请求未进入处理器")
	cancel()

	select {
	case err := <-startErr:
		t.Fatalf("请求完成前服务已退出: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	releaseRequest <- struct{}{}
	require.NoError(t, waitForResult(t, startErr, "服务未在请求完成后退出"))
	require.NoError(t, waitForResult(t, requestDone, "请求未在服务退出后完成"))
}

func TestServer_ShutdownTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())

	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	defer close(releaseRequest)

	srv := New(
		Address(addr),
		Handler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			close(requestStarted)
			<-releaseRequest
		})),
		StopTimeout(50*time.Millisecond),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	startErr := make(chan error, 1)
	go func() {
		startErr <- srv.Start(ctx)
	}()

	waitForServer(t, addr)
	requestDone := requestServer(t, addr)
	waitForSignal(t, requestStarted, "请求未进入处理器")
	cancel()

	assert.ErrorIs(t, waitForResult(t, startErr, "服务未在关闭超时后退出"), context.DeadlineExceeded)
	assert.Error(t, waitForResult(t, requestDone, "请求未在强制关闭后完成"))
}

func requestServer(t *testing.T, addr string) <-chan error {
	t.Helper()

	done := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + addr)
		if err == nil {
			_ = resp.Body.Close()
		}
		done <- err
	}()
	return done
}

func waitForServer(t *testing.T, addr string) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 20*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("服务 %s 未启动", addr)
}

func waitForSignal(t *testing.T, signal <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatal(message)
	}
}

func waitForResult(t *testing.T, result <-chan error, message string) error {
	t.Helper()

	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		t.Fatal(message)
		return nil
	}
}
