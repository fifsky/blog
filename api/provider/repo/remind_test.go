package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestRemind_RemindGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		repo := NewRemind(db)
		ret, err := repo.RemindGetList(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
	})
}
