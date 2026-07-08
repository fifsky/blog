package doubaoasr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientTranscribeSendsExpectedRequest(t *testing.T) {
	t.Parallel()

	var gotHeaders http.Header
	var gotPayload recognizeRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		if r.URL.Path != "/recognize" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("X-Api-Status-Code", "20000000")
		w.Header().Set("X-Api-Message", "OK")
		_, _ = w.Write([]byte(`{"result":{"text":"明天上午九点提醒我喝水。"}}`))
	}))
	defer server.Close()

	client := Client{
		APIKey:     "test-key",
		Endpoint:   server.URL + "/recognize",
		ResourceID: "volc.bigasr.auc_turbo",
		UID:        "tester",
	}

	text, err := client.Transcribe(context.Background(), "YmFzZTY0")
	if err != nil {
		t.Fatalf("Transcribe() error = %v", err)
	}
	if text != "明天上午九点提醒我喝水。" {
		t.Fatalf("unexpected text: %q", text)
	}
	if got := gotHeaders.Get("X-Api-Key"); got != "test-key" {
		t.Fatalf("unexpected X-Api-Key: %q", got)
	}
	if got := gotHeaders.Get("X-Api-Resource-Id"); got != "volc.bigasr.auc_turbo" {
		t.Fatalf("unexpected X-Api-Resource-Id: %q", got)
	}
	if got := gotHeaders.Get("X-Api-Sequence"); got != "-1" {
		t.Fatalf("unexpected X-Api-Sequence: %q", got)
	}
	if gotHeaders.Get("X-Api-Request-Id") == "" {
		t.Fatal("expected X-Api-Request-Id")
	}
	if gotPayload.User.UID != "tester" {
		t.Fatalf("unexpected uid: %q", gotPayload.User.UID)
	}
	if gotPayload.Audio.Data != "YmFzZTY0" {
		t.Fatalf("unexpected audio data: %q", gotPayload.Audio.Data)
	}
	if gotPayload.Request.ModelName != "bigmodel" {
		t.Fatalf("unexpected model name: %q", gotPayload.Request.ModelName)
	}
}

func TestClientTranscribeReturnsHeaderError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Api-Status-Code", "45000001")
		w.Header().Set("X-Api-Message", "invalid request")
		_, _ = w.Write([]byte(`{"message":"invalid"}`))
	}))
	defer server.Close()

	client := Client{APIKey: "test-key", Endpoint: server.URL}
	_, err := client.Transcribe(context.Background(), "YmFzZTY0")
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got == "" || !containsAll(got, "api_status_code=45000001", "api_message=invalid request") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientTranscribeRequiresAPIKey(t *testing.T) {
	t.Parallel()

	client := Client{}
	_, err := client.Transcribe(context.Background(), "YmFzZTY0")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClientTranscribeRejectsEmptyText(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Api-Status-Code", "20000000")
		_, _ = w.Write([]byte(`{"result":{"text":""}}`))
	}))
	defer server.Close()

	client := Client{APIKey: "test-key", Endpoint: server.URL}
	_, err := client.Transcribe(context.Background(), "YmFzZTY0")
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got == "" || !containsAll(got, "empty text") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
