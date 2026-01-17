package doubao

import (
	"app/pkg/jsonutil"
	"context"
	"fmt"
	"os"
	"testing"
)

func TestCreateChatCompletion(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}

	var (
		prompt = `# 角色
根据用户提供的日期查询上海的天气（如：暴雨、雾霾、晚霞）生成一段符合意境的诗句和鼓励的短语，并在最后附上天气信息
1. **信息准确性守护者**：确保提供的信息准确无误。
2. 生成的诗句和短语必须符合意境，不一定要在诗句中包含城市信息，你可以自由发挥。
3. **回答更生动活泼**：请在模型的回复中使用适当的 emoji 标签作为天气和心情的表示 🌟😊🎉，不要在回复中使用格式文本，如**天气信息：**"
`
	)
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
						Text: prompt,
					},
				},
			},
			{
				Role: "user",
				Content: []MessageContent{
					{
						Type: "input_text",
						Text: "日期：2026-01-17",
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
