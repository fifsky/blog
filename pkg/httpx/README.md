# httpx

Simple, composable HTTP client with middleware for Go.

## Features

- Option-based client configuration (timeout, transport, middlewares)
- Composable, chainable middlewares (tracing, mocking, and custom transports)
- Pluggable transport layer (custom TLS, connection tuning)
- Friendly for unit tests via an HTTP mock middleware

## Quick Start

```go
package main

import (
	"net/http"
	"time"

	"app/pkg/httpx"
)

func main() {
	client := httpx.NewClient(
		httpx.WithTimeout(5*time.Second),
		httpx.WithMiddleware(httpx.Trace()),
	)

	req, _ := http.NewRequest(http.MethodGet, "https://httpbin.org/get", nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
```

## Middlewares

The middleware model wraps a `http.RoundTripper` in a chain, allowing cross-cutting features without coupling to request logic.

- HTTP logging: compose any `http.RoundTripper` middleware from caller code
- Trace: OpenTelemetry tracing for HTTP client requests
- Mock: programmable mock responses for tests
- Custom middleware: plug in logging, metrics, retries, or other RoundTripper wrappers

### HTTP Logging

```go
import (
	"bytes"
	"log/slog"
	"net/http"
	"os"

	"app/pkg/httpx"
	"github.com/goapt/logger"
	"github.com/goapt/logger/sloghttp"
)

func exampleAccessLog() {
	l := logger.New(logger.NewJSONHandler(os.Stdout, logger.WithLevel(slog.LevelInfo)))
	client := httpx.NewClient(httpx.WithMiddleware(func(next http.RoundTripper) http.RoundTripper {
		return sloghttp.NewRoundTripper(l, next, sloghttp.DefaultConfig)
	}))

	req, _ := http.NewRequest(http.MethodPost, "https://httpbin.org/anything", bytes.NewReader([]byte(`{"k":"v"}`)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}
```

### Trace (OpenTelemetry)

```go
client := httpx.NewClient(httpx.WithMiddleware(httpx.Trace()))
resp, _ := client.Get("https://httpbin.org/get")
defer resp.Body.Close()
```

### HTTP Mock (for tests)

```go
import (
	"bytes"
	"errors"

	"app/pkg/httpx"
)

func exampleMock() {
	suites := []httpx.MockSuite{
		{URI: "/get", ResponseBody: "ok"},
		{URI: "/user/id/.*", ResponseBody: "user"},
		{URI: "/find\\?id=.*", ResponseBody: "find"},
		{URI: "/bodymatch", MatchBody: map[string]any{"user_id": 1}, ResponseBody: "body-ok"},
		{URI: "/query", MatchQuery: map[string]any{"name": "test"}, ResponseBody: "query-ok"},
		{URI: "/error", Error: errors.New("mock error")},
	}

	client := httpx.NewClient(httpx.WithMiddleware(httpx.Mock(suites)))
	body := bytes.NewBufferString(`{"user_id":1}`)
	resp, _ := client.Post("/bodymatch", "application/json; charset=utf-8", body)
	defer resp.Body.Close()
}
```

## Custom Middleware

```go
import (
	"net/http"

	"app/pkg/httpx"
)

func exampleCustomMW() {
	logMW := func(next http.RoundTripper) http.RoundTripper {
		return httpx.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			resp, err := next.RoundTrip(req)
			return resp, err
		})
	}

	client := httpx.NewClient(httpx.WithMiddleware(logMW))
	_ = client
}
```

## Custom Transport

```go
import (
	"crypto/tls"
	"net/http"

	"app/pkg/httpx"
)

func exampleTransport() {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		CipherSuites: []uint16{tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384},
	}

	client := httpx.NewClient(httpx.WithTransport(tr))
	resp, _ := client.Get("https://httpbin.org/json")
	defer resp.Body.Close()
}
```
