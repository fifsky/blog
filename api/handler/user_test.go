package handler

import (
	"net/http"
	"testing"

	"app/provider/repo"
	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

func TestUser_Login(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema())
		repoUser := repo.NewUser(db)
		handler := NewUser(db, repoUser)

		req := test.NewRequest("/login", handler.Login)
		resp, err := req.JSON(map[string]interface{}{"user_name": "test", "password": "test"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, `{"code":200,"data":{"access_token":"i+LAmF8PLBrVbZcIe88JMpK0coo9wH7yyUNn0z2oxWSmSDA6MqTSksQAZmAQWxok","user":{"id":1,"name":"test","password":"098f6bcd4621d373cade4e832627b4f6","nick_name":"test","email":"test@aaaa.com","status":1,"type":1,"created_at":"2017-08-18 15:21:56","updated_at":"2017-08-23 17:59:58"}},"msg":"success"}`, resp.GetBodyString())
	})
}
