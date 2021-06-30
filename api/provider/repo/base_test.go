package repo

import (
	"testing"
	"time"

	"app/provider/model"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/ilibs/gosql/v2"
	"github.com/stretchr/testify/assert"
)

func TestBase_Create(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repo := &Base{
			db: db,
		}

		user := &model.Users{
			Id:        789,
			Name:      "demo",
			Password:  "123",
			NickName:  "123",
			Email:     "123@123.com",
			Status:    1,
			Type:      1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		id, err := repo.Create(user)
		assert.NoError(t, err)
		assert.Equal(t, int64(user.Id), id)
	})
}

func TestBase_Update(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repo := &Base{
			db: db,
		}

		user := &model.Users{
			Id:   1,
			Name: "demo",
		}

		aff, err := repo.Update(user)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), aff)
	})
}

func TestBase_Find(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repo := &Base{
			db: db,
		}

		user := &model.Users{
			Id: 1,
		}

		err := repo.Find(user)
		assert.NoError(t, err)
		assert.Equal(t, "test", user.Name)
	})
}

func TestBase_FindAll(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repo := &Base{
			db: db,
		}

		user := make([]model.Users, 0)

		err := repo.FindAll(&user, func(b *gosql.Builder) {
			b.Where("status = 1")
		})
		assert.NoError(t, err)
		assert.True(t, len(user) > 0)
	})
}
