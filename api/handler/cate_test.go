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

func TestCate_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		handler := NewCate(db, repo.NewCate(db))
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
				req := test.NewRequest("/api/cate/all", handler.All)
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

func TestCate_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		handler := NewCate(db, repo.NewCate(db))
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
				req := test.NewRequest("/api/admin/cate/list", handler.List)
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

func TestCate_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		handler := NewCate(db, repo.NewCate(db))

		tests := []testutil.TestCase{
			{
				Name:         "success",
				RequestBody:  gee.H{"name": "demo", "domain": "demo", "desc": "demo", "created_at": "2021-06-29 11:55:09", "updated_at": "2021-06-29 11:55:09"},
				ResponseBody: `{"code":200,"data":{"id":2,"name":"demo","desc":"demo","domain":"demo","created_at":"2021-06-29 11:55:09","updated_at":"2021-06-29 11:55:09"},"msg":"success"}`,
			},
			{
				Name:         "update success",
				RequestBody:  gee.H{"id": 1, "domain": "test2", "name": "test2", "desc": "test2", "updated_at": "2021-06-29 11:55:09"},
				AssertType:   testutil.AssertContains,
				ResponseBody: `"name":"test2"`,
			},
			{
				Name:         "update error",
				RequestBody:  gee.H{"id": 1, "domain": "demo", "name": "demo", "desc": "demo"},
				ResponseBody: `{"code":201,"msg":"更新失败"}`,
			},
			{
				Name:         "params error",
				RequestBody:  gee.H{"name": "demo", "desc": "demo"},
				ResponseBody: `{"code":201,"msg":"参数错误:缺少domain"}`,
			},
		}

		for _, tt := range tests {
			tt.Run(t, "/api/admin/cate/post", handler.Post)
		}
	})
}

func TestCate_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates", "posts")...)
		handler := NewCate(db, repo.NewCate(db))
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
			{
				"delete error",
				gee.H{"id": 1},
				`{"code":201,"msg":"该分类下面还有文章，不能删除"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/cate/delete", handler.Delete)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}
