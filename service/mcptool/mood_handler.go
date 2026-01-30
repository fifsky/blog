package mcptool

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"app/store"
	"app/store/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MoodCreateInput struct {
	Content string `json:"content" jsonschema:"要记录的内容，可以是用户说的原话，也可以是提炼后的精华句子"`
}

func NewMoodHandler(s *store.Store) http.Handler {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "blog-mood-mcp",
		Version: "v1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name: "mood_create",
		Description: `记录一条有价值的心情、说说或感悟。返回 JSON：{"id":123}

适合记录的内容类型：
- 情感表达：表达幸福、感恩、快乐、忧伤、思念等情绪的句子
- 人生感悟：对生活、人生、幸福、爱情的思考和领悟
- 哲理金句：富有哲理、发人深省的话语
- 心灵鸡汤：励志、正能量、治愈系的内容
- 此刻心情：用户想要记录的当下感受

判断标准：
1. 内容有情感价值或思想深度
2. 适合作为朋友圈或说说发布
3. 用户表达的是一种感受、态度或见解

注意：不要把普通的问答、任务指令、技术讨论误判为心情`,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input MoodCreateInput) (*mcp.CallToolResult, any, error) {
		if input.Content == "" {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "心情内容不能为空"},
				},
			}, nil, nil
		}

		md := &model.Mood{
			Content:   input.Content,
			UserId:    1, // 默认用户ID
			CreatedAt: time.Now(),
		}
		lastID, err := s.CreateMood(ctx, md)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{
					&mcp.TextContent{Text: "创建心情失败: " + err.Error()},
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
