package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app/config"
	"app/model"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestUser_Login(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		conf := &config.Config{}
		conf.Common.TokenSecret = "abcdabcdabcdabcd"
		handler := NewUser(store.New(db), conf)
		rr := doJSON(handler.Login, "/api/login", map[string]any{"user_name": "test", "password": "test"})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Login, "/api/login", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
		rr3 := doJSON(handler.Login, "/api/login", map[string]any{"user_name": "test", "password": "test234"})
		if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`"code":202`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		}
		rr4 := doJSON(handler.Login, "/api/login", map[string]any{"user_name": "stop", "password": "test"})
		if rr4.Code == http.StatusOK || !bytes.Contains(rr4.Body.Bytes(), []byte(`"code":202`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr4.Code, rr4.Body.String())
		}
	})
}

func TestUser_LoginUser(t *testing.T) {
	user := &model.User{
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
	userHandler := &User{}
	req := httptest.NewRequest(http.MethodGet, "/api/admin/loginuser", nil)
	req = req.WithContext(SetLoginUser(req.Context(), user))
	rr := httptest.NewRecorder()
	userHandler.LoginUser(rr, req)
	if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
		t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestUser_Get(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)

		handler := NewUser(store.New(db), nil)
		rr := doJSON(handler.Get, "/api/admin/user/get", map[string]any{"id": 1})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Get, "/api/admin/user/get", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
		rr3 := doJSON(handler.Get, "/api/admin/user/get", map[string]any{"id": 888})
		if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`"用户不存在"`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		}
	})
}

func TestUser_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(store.New(db), nil)
		rr := doJSON(handler.List, "/api/admin/user/list", map[string]any{"page": 1})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"list"`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestUser_Status(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(store.New(db), nil)
		rr := doJSON(handler.Status, "/api/admin/user/status", map[string]any{"id": 1})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Status, "/api/admin/user/status", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
		rr3 := doJSON(handler.Status, "/api/admin/user/status", map[string]any{"id": 888})
		if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`"code":202`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		}
	})
}

func TestUser_Create(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("users")...)
		handler := NewUser(store.New(db), nil)
		rr := doJSON(handler.Create, "/api/admin/user/create", map[string]any{"name": "demo", "password": "123", "nick_name": "demo", "email": "demo@123.com", "type": 1})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr3 := doJSON(handler.Create, "/api/admin/user/create", map[string]any{"name": "demo", "nick_name": "demo", "type": 1})
		if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`"密码不能为空"`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		}
	})
}
