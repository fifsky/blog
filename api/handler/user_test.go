package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/provider/model"
	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Login(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(db, repo.NewUser(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"user_name": "test", "password": "test"},
				`{"code":200,"data":{"access_token":"i+LAmF8PLBrVbZcIe88JMpK0coo9wH7yyUNn0z2oxWSmSDA6MqTSksQAZmAQWxok","user":{"id":1,"name":"test","password":"098f6bcd4621d373cade4e832627b4f6","nick_name":"test","email":"test@aaaa.com","status":1,"type":1,"created_at":"2017-08-18 15:21:56","updated_at":"2017-08-23 17:59:58"}},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{},
				`{"code":201,"msg":"参数错误:缺少user_name"}`,
			},
			{
				"password error",
				gee.H{"user_name": "test", "password": "test234"},
				`{"code":202,"msg":"用户名或密码错误"}`,
			},
			{
				"delete error",
				gee.H{"user_name": "stop", "password": "test"},
				`{"code":202,"msg":"用户已停用"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/login", handler.Login)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestUser_LoginUser(t *testing.T) {
	ctx, _ := gee.CreateTestContext(httptest.NewRecorder())

	user := &model.Users{
		Id:        1,
		Name:      "test",
		Password:  "test",
		NickName:  "test",
		Email:     "test@test.com",
		Status:    1,
		Type:      1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	ctx.Set("userInfo", user)

	userHandler := &User{}
	resp := userHandler.LoginUser(ctx)
	resp.Render()
	assert.Equal(t, `{"code":200,"data":{"id":1,"name":"test","password":"test","nick_name":"test","email":"test@test.com","status":1,"type":1,"created_at":"0001-01-01 08:05:43","updated_at":"0001-01-01 08:05:43"},"msg":"success"}`, string(ctx.ResponseBody()))
}

func TestUser_Get(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(db, repo.NewUser(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"id": 1},
				`{"code":200,"data":{"id":1,"name":"test","password":"098f6bcd4621d373cade4e832627b4f6","nick_name":"test","email":"test@aaaa.com","status":1,"type":1,"created_at":"2017-08-18 15:21:56","updated_at":"2017-08-23 17:59:58"},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{},
				`{"code":201,"msg":"参数错误:缺少id"}`,
			},
			{
				"user not found",
				gee.H{"id": 888},
				`{"code":201,"msg":"用户不存在"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/user/get", handler.Get)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestUser_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(db, repo.NewUser(db))
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
				req := test.NewRequest("/api/admin/user/list", handler.List)
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

func TestUser_Status(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(db, repo.NewUser(db))
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
			{
				"user not found",
				gee.H{"id": 888},
				`{"code":202,"msg":"用户不存在"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/user/status", handler.Status)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}

func TestUser_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(db, repo.NewUser(db))
		tests := []struct {
			name         string
			requestBody  gee.H
			responseBody string
		}{
			{
				"success",
				gee.H{"name": "demo", "password": "123", "nick_name": "demo", "email": "demo@123.com", "type": 1, "created_at": "2021-06-29 11:55:09", "updated_at": "2021-06-29 11:55:09"},
				`{"code":200,"data":{"id":4,"name":"demo","password":"202cb962ac59075b964b07152d234b70","nick_name":"demo","email":"demo@123.com","status":0,"type":1,"created_at":"2021-06-29 11:55:09","updated_at":"2021-06-29 11:55:09"},"msg":"success"}`,
			},
			{
				"params error",
				gee.H{"name": "demo2", "password": "123", "type": 1},
				`{"code":201,"msg":"参数错误:缺少nick_name"}`,
			},
			{
				"password error",
				gee.H{"name": "demo", "nick_name": "demo", "type": 1},
				`{"code":201,"msg":"密码不能为空"}`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := test.NewRequest("/api/admin/user/post", handler.Post)
				resp, err := req.JSON(tt.requestBody)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, resp.Code)
				require.Equal(t, tt.responseBody, resp.GetBodyString())
			})
		}
	})
}
