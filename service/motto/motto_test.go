package motto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"app/pkg/agent"
	"app/pkg/dbunit"
	"app/testutil"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIProvider_GenerateWrapsAgentRun(t *testing.T) {
	type chatRequest struct {
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}

	var requests []chatRequest
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		requests = append(requests, req)

		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"Test Motto\"},\"finish_reason\":null}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer ts.Close()

	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options"))
		_, err := db.ExecContext(context.Background(), `insert into options (id, option_key, option_value) values
			(20, 'ai_token', 'test-token'),
			(21, 'ai_endpoint', ?),
			(22, 'ai_model', 'test-model')`, ts.URL)
		require.NoError(t, err)

		provider := NewOpenAIProvider(agent.New(
			agent.WithClient(openai.NewClient(option.WithAPIKey("test"), option.WithBaseURL(ts.URL))),
			agent.WithModel("test"),
		))
		got, err := provider.Generate(context.Background(), "system prompt", "2026-05-22")
		require.NoError(t, err)
		assert.Equal(t, "Test Motto", got)

		require.Len(t, requests, 1)
		assert.Equal(t, "system", requests[0].Messages[0].Role)
		assert.Equal(t, "user", requests[0].Messages[1].Role)
		assert.Equal(t, "2026-05-22", requests[0].Messages[1].Content)
	})
}

func TestOpenAIProvider_Generate(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test")
	}
	if os.Getenv("AI_TOKEN") == "" {
		t.Skip("skip test due to missing AI_TOKEN")
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("AI_TOKEN")),
		option.WithBaseURL(os.Getenv("AI_ENDPOINT")),
	)
	ai := NewOpenAIProvider(agent.New(
		agent.WithClient(client),
		agent.WithModel(os.Getenv("AI_MODEL")),
		agent.WithMCP(map[string]agent.MCPConfig{
			"web_search": {
				Name: "联网搜索",
				URL:  os.Getenv("WEBSEARCH_MCP"),
			},
		}),
	))

	content, err := ai.Generate(context.Background(), Prompt, time.Now().Format("2006-01-02"))
	require.NoError(t, err)
	fmt.Println(content)
}
