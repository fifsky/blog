package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCorsPreflight(t *testing.T) {
	_ = os.Setenv("APP_ENV", "prod")
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := NewCors(next)

	req := httptest.NewRequest(http.MethodOptions, "https://api.fifsky.com/api/mood/list", nil)
	req.Header.Set("Origin", "https://fifsky.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://fifsky.com" {
		t.Fatalf("expected A-C-Allow-Origin https://fifsky.com, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET" {
		t.Fatalf("expected A-C-Allow-Methods GET, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got != "Authorization,Content-Type" {
		t.Fatalf("expected A-C-Allow-Headers Authorization,Content-Type, got %q", got)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("expected Vary Origin, got %q", got)
	}
}
