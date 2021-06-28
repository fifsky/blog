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
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repoUser := repo.NewUser(db)
		handler := NewUser(db, repoUser)

		req := test.NewRequest("/api/login", handler.Login)
		resp, err := req.JSON(map[string]interface{}{"user_name": "test", "password": "test"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, `{"code":200,"data":{"access_token":"i+LAmF8PLBrVbZcIe88JMpK0coo9wH7yyUNn0z2oxWSmSDA6MqTSksQAZmAQWxok","user":{"id":1,"name":"test","password":"098f6bcd4621d373cade4e832627b4f6","nick_name":"test","email":"test@aaaa.com","status":1,"type":1,"created_at":"2017-08-18 15:21:56","updated_at":"2017-08-23 17:59:58"}},"msg":"success"}`, resp.GetBodyString())
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
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repoUser := repo.NewUser(db)
		handler := NewUser(db, repoUser)
		req := test.NewRequest("/api/admin/user/get", handler.Get)
		resp, err := req.JSON(map[string]interface{}{"id": 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, `{"code":200,"data":{"id":1,"name":"test","password":"098f6bcd4621d373cade4e832627b4f6","nick_name":"test","email":"test@aaaa.com","status":1,"type":1,"created_at":"2017-08-18 15:21:56","updated_at":"2017-08-23 17:59:58"},"msg":"success"}`, resp.GetBodyString())
	})
}

func TestUser_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repoUser := repo.NewUser(db)
		handler := NewUser(db, repoUser)
		req := test.NewRequest("/api/admin/user/list", handler.List)
		resp, err := req.JSON(map[string]interface{}{"page": 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.True(t, resp.GetJsonBody("data.list").Exists())
	})
}

func TestUser_Status(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("users"))
		repoUser := repo.NewUser(db)
		handler := NewUser(db, repoUser)
		req := test.NewRequest("/api/admin/user/status", handler.Status)
		resp, err := req.JSON(map[string]interface{}{"id": 1})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, int64(200), resp.GetJsonBody("code").Int())
	})
}
