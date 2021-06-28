package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestUser_GetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repoUser := NewUser(db)
		users, err := repoUser.GetList(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, "rita", users[0].Name)
	})
}

func TestUser_GetUser(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repoUser := NewUser(db)
		ret, err := repoUser.GetUser(1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, "test", ret.Name)
	})
}
