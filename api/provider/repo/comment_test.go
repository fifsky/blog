package repo

import (
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/stretchr/testify/assert"
)

func TestComment_PostComments(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments")...)
		repo := NewComment(db)
		ret, err := repo.PostComments(7, 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestComment_NewComments(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "users", "posts")...)
		db.ShowSql()
		repo := NewComment(db)
		ret, err := repo.NewComments()
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestComment_CommentList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "users", "posts")...)
		db.ShowSql()
		repo := NewComment(db)
		ret, err := repo.CommentList(1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestComment_PostCommentNum(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "users", "posts")...)
		db.ShowSql()
		repo := NewComment(db)
		ret, err := repo.PostCommentNum([]int{7})
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}
