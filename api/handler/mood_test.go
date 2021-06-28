package handler

import (
	"net/http"
	"testing"

	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

func TestMood_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users")...)
		handler := NewMood(db, repo.NewMood(db))
		req := test.NewRequest("/api/mood/list", handler.List)
		resp, err := req.JSON(gee.H{"page": 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, int64(200), resp.GetJsonBody("code").Int())
		require.True(t, resp.GetJsonBody("data.list").Exists())
	})
}
