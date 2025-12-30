package store

import (
	"context"
	"testing"

	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestMood_MoodGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users")...)
		s := New(db)
		ret, err := s.ListMood(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}
