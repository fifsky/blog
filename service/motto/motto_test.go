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

	"app/config"
	"app/pkg/bark"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAIProvider
type MockAIProvider struct {
	Result string
	Err    error
}

func (m *MockAIProvider) Generate(ctx context.Context, prompt, content string) (string, error) {
	return m.Result, m.Err
}

func TestMotto_GenerateDailyMotto(t *testing.T) {
	// Mock Bark Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Prepare DB
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users")...)
		s := store.New(db)

		// Prepare Config
		conf := &config.Config{}

		// Prepare Bark Client
		barkClient := bark.New(http.DefaultClient, ts.URL, "test-token")

		// Prepare Mock AI
		ai := &MockAIProvider{
			Result: "Test Motto Content",
		}

		m := New(s, conf, barkClient, ai)

		// Execute
		err := m.GenerateDailyMotto()
		assert.NoError(t, err)

		// Verify DB
		moods, err := s.ListMood(context.Background(), 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, moods)
		assert.Equal(t, "Test Motto Content", moods[0].Content)
	})
}

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
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		_, err := db.ExecContext(context.Background(), `insert into options (id, option_key, option_value) values
			(20, 'ai_token', 'test-token'),
			(21, 'ai_endpoint', ?),
			(22, 'ai_model', 'test-model')`, ts.URL)
		require.NoError(t, err)

		provider := NewOpenAIProvider(&config.Config{}, store.New(db))
		got, err := provider.Generate(context.Background(), "system prompt", "2026-05-22")
		require.NoError(t, err)
		assert.Equal(t, "Test Motto", got)

		require.Len(t, provider.history, 2)
		require.Len(t, requests, 1)
		assert.Equal(t, "system", requests[0].Messages[0].Role)
		assert.Equal(t, "user", requests[0].Messages[1].Role)
		assert.Equal(t, "2026-05-22", requests[0].Messages[1].Content)
	})
}

func TestOpenAIProvider_Generate(t *testing.T) {
	// t.Skip("skip test")
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("AI_TOKEN")),
		option.WithBaseURL(os.Getenv("AI_ENDPOINT")),
	)
	req := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(time.Now().Format(time.DateOnly)),
		},
		Model:           openai.ChatModel(os.Getenv("AI_MODEL")),
		ReasoningEffort: "high",
	}

	req.SetExtraFields(map[string]any{
		"thinking": map[string]any{
			"type": "enabled",
		},
	})

	content, err := client.Chat.Completions.New(context.Background(), req)
	require.NoError(t, err)
	fmt.Println(content.Choices[0].Message.Content)
}
