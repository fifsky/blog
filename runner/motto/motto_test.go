package motto

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/store"
	"app/testutil"

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
	// Prepare DB
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users"))
		s := store.New(db)

		// Prepare Mock AI
		ai := &MockAIProvider{
			Result: "Test Motto Content",
		}

		m := New(s, ai, "0 7 * * *")

		// Execute
		err := m.generateDailyMotto()
		assert.NoError(t, err)

		// Verify DB
		moods, err := s.ListMood(context.Background(), 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, moods)
		assert.Equal(t, "Test Motto Content", moods[0].Content)
	})
}
