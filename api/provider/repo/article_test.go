package repo

import (
	"testing"

	"app/provider/model"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestArticle_GetUserPost(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts", "users", "cates")...)
		commRepo := NewComment(db)
		repo := NewArticle(db, commRepo)
		ret, err := repo.GetUserPost(4, "")
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.NotNil(t, ret.User)
		assert.NotNil(t, ret.Cate)
		assert.Equal(t, ret.CateId, 1)
		assert.Equal(t, ret.UserId, 1)
	})
}

func TestArticle_PostPrev(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts")...)
		commRepo := NewComment(db)
		repo := NewArticle(db, commRepo)
		ret, err := repo.PostPrev(7)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret.Id, 4)
	})
}

func TestArticle_PostNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts")...)
		commRepo := NewComment(db)
		repo := NewArticle(db, commRepo)
		ret, err := repo.PostNext(4)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret.Id, 8)
	})
}

func TestArticle_PostArchive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts")...)
		commRepo := NewComment(db)
		repo := NewArticle(db, commRepo)
		ret, err := repo.PostArchive()
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestArticle_PostGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts", "users", "cates")...)
		commRepo := NewComment(db)
		repo := NewArticle(db, commRepo)

		p := &model.Posts{
			CateId: 1,
		}

		ret, err := repo.PostGetList(p, 1, 1, "", "")
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}
