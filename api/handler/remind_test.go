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

func setRemind(c *gee.Context) gee.Response {
	c.Set("remind", &model.Reminds{
		Id:        1,
		Type:      1,
		Content:   "demo",
		Month:     1,
		Week:      0,
		Day:       1,
		Hour:      1,
		Minute:    1,
		Status:    2,
		CreatedAt: timeutil.MustParseDateTime("2020-10-24 12:23:34"),
	})
	return nil
}

func TestRemind_Change(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(db, repo.NewRemind(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{},
				`已确认收到提醒`,
			},
			{
				"not found",
				gee.H{},
				`{"code":202,"msg":"记录未找到"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var req *test.Request
				if tt.name == "not found" {
					req = test.NewRequest("/api/remind/change", handler.Change)
				} else {
					req = test.NewRequest("/api/remind/change", setRemind, handler.Change)
				}
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestRemind_Delay(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(db, repo.NewRemind(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{},
				`将在10分钟后再次提醒`,
			},
			{
				"not found",
				gee.H{},
				`{"code":202,"msg":"记录未找到"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var req *test.Request
				if tt.name == "not found" {
					req = test.NewRequest("/api/remind/delay", handler.Delay)
				} else {
					req = test.NewRequest("/api/remind/delay", setRemind, handler.Delay)
				}
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestRemind_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(db, repo.NewRemind(db))
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
				req := test.NewRequest("/api/admin/remind/list", handler.List)
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

func TestRemind_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(db, repo.NewRemind(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"type": 1, "content": "demo", "month": 1, "week": 0, "day": 1, "hour": 1, "minute": 1, "status": 0, "created_at": "2021-06-29 11:55:09"},
				`{"code":200,"data":{"id":10,"type":1,"content":"demo","month":1,"week":0,"day":1,"hour":1,"minute":1,"status":0,"next_time":"0001-01-01T08:05:43+08:05","created_at":"2021-06-29 11:55:09"},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{"content": "demo"},
				`{"code":201,"msg":"参数错误:缺少type"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/remind/post", handler.Post)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestRemind_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("reminds"))
		handler := NewRemind(db, repo.NewRemind(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"id": 8},
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
				req := test.NewRequest("/api/admin/remind/delete", handler.Delete)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}
