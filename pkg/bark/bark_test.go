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

	clientb := New(http.DefaultClient, os.Getenv("BARK_URL"), os.Getenv("BARK_TOKEN"))
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
