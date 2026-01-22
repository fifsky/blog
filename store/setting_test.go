package store

import (
	"context"
	"testing"

	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestSetting_GetOptions(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixtures("options")...)
		s := New(db)
		ret, err := s.GetOptions(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret["site_name"], "無處告別")
	})
}
