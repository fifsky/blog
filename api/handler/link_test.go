package handler

import (
	"net/http"
	"testing"

	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLink_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(db, repo.NewLink(db))
		tests := []struct {
			name        string
			requestBody gee.H
			checkFunc   func(t *testing.T, resp *test.Response)
		}{
			{
				"success",
				gee.H{},
				func(t *testing.T, resp *test.Response) {
					assert.True(t, len(resp.GetJsonBody("data").Array()) > 0)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/link/all", handler.All)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				if tt.checkFunc != nil {
					tt.checkFunc(t, resp)
				}
			})
		}
	})
}

func TestLink_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(db, repo.NewLink(db))
		tests := []struct {
			name        string
			requestBody gee.H
			checkFunc   func(t *testing.T, resp *test.Response)
		}{
			{
				"success",
				gee.H{},
				func(t *testing.T, resp *test.Response) {
					assert.True(t, len(resp.GetJsonBody("data.list").Array()) > 0)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/link/list", handler.List)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				if tt.checkFunc != nil {
					tt.checkFunc(t, resp)
				}
			})
		}
	})
}

func TestLink_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(db, repo.NewLink(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"name": "demo", "url": "https://example.com", "desc": "demo", "created_at": "2021-06-29 11:55:09", "updated_at": "2021-06-29 11:55:09"},
				`{"code":200,"data":{"id":4,"name":"demo","url":"https://example.com","desc":"demo","created_at":"2021-06-29 11:55:09"},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{"name": "demo", "desc": "demo"},
				`{"code":201,"msg":"参数错误:缺少url"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/link/post", handler.Post)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestLink_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(db, repo.NewLink(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"id": 3},
				`{"code":200,"msg":"success"}`,
			},
			{
				"params error",
				gee.H{},
				`{"code":201,"msg":"参数错误:缺少id"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/link/delete", handler.Delete)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}
