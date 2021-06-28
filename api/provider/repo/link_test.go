package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestLink_GetAllLinks(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		repo := NewLink(db)
		ret, err := repo.GetAllLinks()
		assert.NoError(t, err)
		assert.NotNil(t, ret)
	})
}
