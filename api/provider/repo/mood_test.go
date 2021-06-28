package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestMood_MoodGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users")...)
		repo := NewMood(db)
		ret, err := repo.MoodGetList(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.NotNil(t, ret[0].User)
	})
}
