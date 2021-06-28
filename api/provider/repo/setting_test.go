package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestSetting_GetOptions(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		repo := NewSetting(db)
		ret, err := repo.GetOptions()
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret["site_name"], "無處告別")
	})
}
