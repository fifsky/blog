package httpx_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"app/pkg/httpx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestHttpError(t *testing.T) {
	client := httpx.NewClient(httpx.WithTimeout(time.Second * 1))

	got, err := client.Post("https://1234567890000.com/test", "application/json; charset=utf-8", bytes.NewReader([]byte(`{"name":"test"}`)))
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestNewClientWithTarce(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external network test in short mode")
	}
	sr := &tracetest.SpanRecorder{}
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	var tr = otel.Tracer("httplib/client")
	ctx, span := tr.Start(context.Background(), "test.external")
	span.SetAttributes(
		attribute.String("test.message", "hello"),
	)
	span.End()

	client := httpx.NewClient(httpx.WithMiddleware(httpx.Trace()))
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
	if testing.Short() {
		t.Skip("skipping external network test in short mode")
	}
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

	client := httpx.NewClient(httpx.WithTransport(transport))

	got, err := client.Get("https://httpbin.org/json")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, got.StatusCode)
}
