package handler

import (
	"net/http"
	"strings"
	"testing"

	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticle_Archive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
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
				req := test.NewRequest("/api/article/archive", handler.Archive)
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

func TestArticle_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
		tests := []struct {
			name        string
			requestBody gee.H
			checkFunc   func(t *testing.T, resp *test.Response)
		}{
			{
				"success",
				gee.H{"page": 1},
				func(t *testing.T, resp *test.Response) {
					assert.True(t, len(resp.GetJsonBody("data.list").Array()) > 0)
				},
			},
			{
				"params error",
				gee.H{},
				func(t *testing.T, resp *test.Response) {
					assert.Equal(t, `{"code":201,"msg":"参数错误:缺少page"}`, resp.GetBodyString())
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/article/list", handler.List)
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

func TestArticle_PrevNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"id": 7},
				`{"code":200,"data":{"next":{"id":8,"title":"example"},"prev":{"id":4,"title":"fifsky blog for php!"}},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{},
				`{"code":201,"msg":"参数错误:缺少id"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/article/prevnext", handler.PrevNext)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestArticle_Detail(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
		tests := []struct {
			name        string
			requestBody gee.H
			checkFunc   func(t *testing.T, resp *test.Response)
		}{
			{
				"success",
				gee.H{"id": 7},
				func(t *testing.T, resp *test.Response) {
					assert.Equal(t, `关于`, resp.GetJsonBody("data.title").String())
				},
			},
			{
				"params error",
				gee.H{},
				func(t *testing.T, resp *test.Response) {
					assert.Equal(t, `{"code":201,"msg":"参数错误"}`, resp.GetBodyString())
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/article/detail", handler.Detail)
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

func TestArticle_Feed(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts", "users")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
		tests := []struct {
			name        string
			requestBody gee.H
			checkFunc   func(t *testing.T, resp *test.Response)
		}{
			{
				"success",
				gee.H{},
				func(t *testing.T, resp *test.Response) {
					assert.True(t, strings.Contains(resp.GetBodyString(), "<title>关于</title>"))
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/feed.xml", handler.Feed)
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

func TestArticle_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"cate_id": 1, "type": 1, "title": "test", "url": "", "content": "test", "created_at": "2021-06-29 11:55:09", "updated_at": "2021-06-29 11:55:09"},
				`{"code":200,"data":{"id":9,"cate_id":1,"type":1,"user_id":1,"title":"test","url":"","content":"test","status":1,"created_at":"2021-06-29 11:55:09","updated_at":"2021-06-29 11:55:09"},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{"cate_id": 1, "type": 1, "title": "", "url": "", "content": "test"},
				`{"code":201,"msg":"参数错误:缺少title"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/article/post", setLoginUser, handler.Post)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestArticle_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options", "posts")...)
		handler := NewArticle(db, repo.NewArticle(db, repo.NewComment(db)), repo.NewSetting(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"id": 4},
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
				req := test.NewRequest("/api/admin/article/delete", handler.Delete)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}
