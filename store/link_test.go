package store

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLink_GetAllLinks(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		s := New(db)
		ret, err := s.GetAllLinks(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
	})
}
