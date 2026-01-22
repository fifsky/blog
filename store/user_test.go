package store

import (
	"context"
	"testing"

	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestUser_GetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixture("users"))
		s := New(db)
		users, err := s.ListUser(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.True(t, len(users) > 0)
	})
}

func TestUser_GetUser(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.TestDSN, testutil.Schema(), testutil.Fixture("users"))
		s := New(db)
		ret, err := s.GetUser(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, "test", ret.Name)
	})
}
