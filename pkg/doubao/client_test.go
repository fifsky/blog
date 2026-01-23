package doubao

import (
	"app/pkg/jsonutil"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateChatCompletion(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}
	if os.Getenv("AI_TOKEN") == "" {
		t.Skip("skip integration test: AI_TOKEN not set")
	}

	prompt, err := os.ReadFile("testdata/blog.md")
	require.NoError(t, err)
	userInput := "请搜索2026年1月18日关于Golang的技术文章，并生成一篇1000字左右的技术博文"
	// Create client with mock server URL
	client := NewClient(os.Getenv("AI_TOKEN"))
	// Request data matches the curl example structure
	req := &ChatRequest{
		Model: "doubao-seed-1-8-251228",
		Tools: []Tool{
			{
				Type:       "web_search",
				MaxKeyword: 2,
				Limit:      2,
			},
		},
		MaxToolCalls: 1,
		Thinking: &Thinking{
			Type: "disabled",
		},
		Input: []Message{
			{
				Role: "system",
				Content: []MessageContent{
					{
						Type: "input_text",
						Text: string(prompt),
					},
				},
			},
			{
				Role: "user",
				Content: []MessageContent{
					{
						Type: "input_text",
						Text: userInput,
					},
				},
			},
		},
	}

	// Call API
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateChatCompletion failed: %v", err)
	}
	// Check response
	if len(resp.Output) == 0 {
		t.Fatalf("Expected non-empty output, got empty")
	}
	fmt.Println(jsonutil.Encode(resp))
}
