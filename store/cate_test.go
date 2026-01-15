package store

import (
	"context"
	"testing"

	"app/pkg/jsonutil"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestCate_GetAllCates(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates", "posts")...)
		s := New(db)
		ret, err := s.GetAllCates(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		t.Log(jsonutil.Encode(ret))
	})
}
