package httputil_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"app/pkg/httputil"
	"app/pkg/httputil/middleware"
	"app/pkg/logger"
	"app/pkg/logger/sloghttp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestHttpError(t *testing.T) {
	client := httputil.NewClient(httputil.WithTimeout(time.Second * 1))

	got, err := client.Post("https://1234567890000.com/test", "application/json; charset=utf-8", bytes.NewReader([]byte(`{"name":"test"}`)))
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestRequestWithLogInfo(t *testing.T) {
	buf := strings.Builder{}
	l := logger.New(&logger.Config{
		Mode:   logger.ModeCustom,
		Writer: &buf,
	})

	request := map[string]any{"app_key": "ARkBS09IS0saFExLThQ", "sign": "test", "time": "test"}
	client := httputil.NewClient(httputil.WithMiddleware(middleware.AccessLog(l)))
	b, _ := json.Marshal(request)
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://httpbin.org/anything", bytes.NewReader(b))
	sloghttp.AddCustomAttributes(req, slog.String("other", "is other info"))

	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	assert.NoError(t, err)

	fmt.Println(buf.String())

	m := make(map[string]any)
	err = json.Unmarshal([]byte(buf.String()), &m)
	assert.NoError(t, err)
	assert.Equal(t, "is other info", m["other"])
}

func TestNewClientWithTarce(t *testing.T) {
	sr := &tracetest.SpanRecorder{}
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	var tr = otel.Tracer("httplib/client")
	ctx, span := tr.Start(context.Background(), "test.external")
	span.SetAttributes(
		attribute.String("test.message", "hello"),
	)
	span.End()

	client := httputil.NewClient(httputil.WithMiddleware(middleware.Debug(), middleware.Trace()))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://httpbin.org/json", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	_, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	spans := sr.Ended()

	for _, ss := range spans {
		b, _ := json.Marshal(ss.Attributes())
		fmt.Println(ss.Name(), string(b))
	}
}

func TestWithTransport(t *testing.T) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		CipherSuites: []uint16{
			// TLS 1.0 - 1.2 cipher suites.
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			// TLS 1.3 cipher suites.
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	client := httputil.NewClient(httputil.WithTransport(transport))

	got, err := client.Get("https://httpbin.org/json")
	assert.NoError(t, err)
	body, err := io.ReadAll(got.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(body))
}
