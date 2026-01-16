package motto

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/config"
	"app/pkg/bark"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
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
