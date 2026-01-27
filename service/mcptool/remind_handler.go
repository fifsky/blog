package mcptool

import (
	"context"
	"encoding/json"
	"net/http"

	"app/config"
	"app/store"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type RemindSmartCreateInput struct {
	Content string `json:"content" jsonschema:"提醒内容，自然语言描述即可"`
}

func NewRemindHandler(conf *config.Config, s *store.Store) http.Handler {
	aiClient := openai.NewClient(
		option.WithAPIKey(conf.Common.AIToken),
		option.WithBaseURL(conf.Common.AIEndpoint),
	)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "blog-remind-mcp",
		Version: "v1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remind_smart_create",
		Description: "根据自然语言创建提醒，返回 JSON：{\"id\":123}",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input RemindSmartCreateInput) (*mcp.CallToolResult, any, error) {
		lastID, err := SmartCreateRemind(ctx, conf, aiClient, s, input.Content)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Reason + ": " + err.Message},
				},
			}, nil, nil
		}

		out, _ := json.Marshal(map[string]any{"id": lastID})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(out)},
			},
		}, nil, nil
	})

	return mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{Stateless: true})
}
