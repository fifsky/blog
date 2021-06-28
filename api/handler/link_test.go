package handler

import (
	"database/sql"
	"net/http"
	"testing"

	"app/provider/model"
	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

func TestLink_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("links"))
		repoLink := repo.NewLink(db)
		handler := NewLink(db, repoLink)
		req := test.NewRequest("/api/link/all", handler.All)
		resp, err := req.JSON(nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, int64(200), resp.GetJsonBody("code").Int())
		require.Equal(t, "圆子", resp.GetJsonBody("data.0.content").String())
	})
}

func TestLink_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("links"))
		repoLink := repo.NewLink(db)
		handler := NewLink(db, repoLink)
		req := test.NewRequest("/api/link/list", handler.List)
		resp, err := req.JSON(map[string]interface{}{"page": 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, int64(200), resp.GetJsonBody("code").Int())
		require.True(t, resp.GetJsonBody("data.list").Exists())
	})
}

func TestLink_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("links"))
		repoLink := repo.NewLink(db)
		handler := NewLink(db, repoLink)
		req := test.NewRequest("/api/admin/link/post", handler.Post)

		link := map[string]interface{}{
			"name": "new link",
			"url":  "https://example.com",
			"desc": "new link demo",
		}

		resp, err := req.JSON(link)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, int64(200), resp.GetJsonBody("code").Int())

		link2 := &model.Links{
			Id: int(resp.GetJsonBody("data.id").Int()),
		}

		err = repoLink.Find(link2)
		require.NoError(t, err)
		require.Equal(t, link["name"], link2.Name)
	})
}

func TestLink_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("links"))
		repoLink := repo.NewLink(db)
		handler := NewLink(db, repoLink)
		req := test.NewRequest("/api/admin/link/delete", handler.Delete)

		resp, err := req.JSON(gee.H{"id": 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, int64(200), resp.GetJsonBody("code").Int())

		link2 := &model.Links{
			Id: 1,
		}

		err = repoLink.Find(link2)
		require.Equal(t, sql.ErrNoRows, err)
	})
}
