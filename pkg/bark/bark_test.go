package bark

import (
	"net/http"
	"os"
	"testing"
)

func TestClient_Send(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	url := os.Getenv("BARK_URL")
	token := os.Getenv("BARK_TOKEN")
	if url == "" || token == "" {
		t.Skip("skip integration test: BARK_URL/BARK_TOKEN not set")
	}
	clientb := New(http.DefaultClient, url, token)
	msg := Message{
		Title:    "test",
		Body:     "test",
		Badge:    1,
		Markdown: "这是测试的",
		Group:    "test",
		Level:    "timeSensitive",
	}
	if err := clientb.Send(msg); err != nil {
		t.Errorf("Client.Send() error = %v", err)
	}
}
