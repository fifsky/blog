package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestCate_GetAllCates(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		repo := NewCate(db)
		ret, err := repo.GetAllCates()
		assert.NoError(t, err)
		assert.NotNil(t, ret)
	})
}
