package motto

import (
	"app/config"
	"app/pkg/bark"
	"app/store"
	"app/testutil"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

func TestOpenAIProvider_Generate(t *testing.T) {
	t.Skip("skip test")
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("AI_TOKEN")),
		option.WithBaseURL(os.Getenv("AI_ENDPOINT")),
	)
	prompt := "每天自动根据用户所在城市的天气（如：暴雨、雾霾、晚霞）生成一段符合意境的诗句或短评。示例： “今日上海大雨。AI 检测到 80% 的忧郁湿度，建议配一杯热可可和坂本龙一的钢琴曲。”"

	req := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage("Test Content"),
		},
		Model: openai.ChatModel(os.Getenv("AI_MODEL")),
		Tools: []openai.ChatCompletionToolUnionParam{},
	}

	req.SetExtraFields(map[string]any{
		"thinking": map[string]any{
			"type": "disabled",
		},
	})

	content, err := client.Chat.Completions.New(context.Background(), req)
	require.NoError(t, err)
	fmt.Println(content.Choices[0].Message.Content)
}
