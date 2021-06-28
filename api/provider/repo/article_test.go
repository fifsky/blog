package repo

import (
	"testing"

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
