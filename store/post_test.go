package store

import (
	"context"
	"testing"

	"app/store/model"
	"app/testutil"

	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestArticle_PostPrev(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts")...)
		s := New(db)
		ret, err := s.PrevPost(context.Background(), 7)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret.Id, 4)
	})
}

func TestArticle_PostNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts")...)
		s := New(db)
		ret, err := s.NextPost(context.Background(), 4)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret.Id, 8)
	})
}

func TestArticle_PostArchive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts")...)
		s := New(db)
		ret, err := s.PostArchive(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestArticle_PostGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts", "users", "cates")...)
		s := New(db)

		p := &model.Post{
			CateId: 1,
		}

		ret, err := s.ListPost(context.Background(), p, 1, 1, "", "", "")
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}
