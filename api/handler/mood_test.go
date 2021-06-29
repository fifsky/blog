package handler

import (
	"net/http"
	"testing"

	"app/provider/model"
	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/golib/timeutil"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setLoginUser(c *gee.Context) gee.Response {
	c.Set("userInfo", &model.Users{
		Id:        1,
		Name:      "test",
		Password:  "123123",
		NickName:  "test",
		Email:     "test@gmail",
		Status:    1,
		Type:      1,
		CreatedAt: timeutil.MustParseDateTime("2020-10-24 12:23:34"),
		UpdatedAt: timeutil.MustParseDateTime("2020-10-24 12:23:34"),
	})
	return nil
}

func TestMood_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users")...)
		handler := NewMood(db, repo.NewMood(db))
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
				req := test.NewRequest("/api/admin/mood/list", setLoginUser, handler.List)
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

func TestMood_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("moods"))
		handler := NewMood(db, repo.NewMood(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"content": "demo"},
				`{"code":200,"msg":"success"}`,
			},
			{
				"params error",
				gee.H{},
				`{"code":201,"msg":"参数错误:缺少content"}`,
			},
			{
				"update",
				gee.H{"id": 1, "content": "demo2"},
				`{"code":200,"msg":"success"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/mood/post", setLoginUser, handler.Post)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestMood_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("moods"))
		handler := NewMood(db, repo.NewMood(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"id": 1},
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
				req := test.NewRequest("/api/admin/mood/delete", setLoginUser, handler.Delete)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}
