package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/model"
	"app/pkg/wechat"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestRemind_Change(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		robot := wechat.NewRobot("123")
		handler := NewRemind(store.New(db), robot)
		// success with remind context
		req := httptest.NewRequest(http.MethodPost, "/api/remind/change", bytes.NewReader([]byte(`{}`)))
		req = req.WithContext(context.WithValue(req.Context(), "remind", &model.Remind{Id: 1}))
		rr := httptest.NewRecorder()
		handler.Change(rr, req)
		if rr.Code != http.StatusOK || rr.Body.String() == "" {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		// not found
		rr2 := doJSON(handler.Change, "/api/remind/change", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`记录未找到`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
	})
}

func TestRemind_Delay(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		robot := wechat.NewRobot("123")
		handler := NewRemind(store.New(db), robot)
		rq := httptest.NewRequest(http.MethodPost, "/api/remind/delay", bytes.NewReader([]byte(`{}`)))
		rq = rq.WithContext(context.WithValue(rq.Context(), "remind", &model.Remind{Id: 1}))
		rr := httptest.NewRecorder()
		handler.Delay(rr, rq)
		if rr.Code != http.StatusOK || rr.Body.String() == "" {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Delay, "/api/remind/delay", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`记录未找到`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
	})
}

func TestRemind_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(store.New(db), wechat.NewRobot("123"))
		rr := doJSON(handler.List, "/api/admin/remind/list", map[string]any{"page": 1})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"list"`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestRemind_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(store.New(db), wechat.NewRobot("123"))
		rr := doJSON(handler.Post, "/api/admin/remind/post", map[string]any{"type": 1, "content": "demo", "month": 1, "week": 0, "day": 1, "hour": 1, "minute": 1, "status": 0, "created_at": "2021-06-29 11:55:09"})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		// rr2 := doJSON(handler.Post, "/api/admin/remind/post", map[string]any{"type": 1})
		// if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
		// 	t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		// }
	})
}

func TestRemind_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("reminds"))
		handler := NewRemind(store.New(db), wechat.NewRobot("123"))
		rr := doJSON(handler.Delete, "/api/admin/remind/delete", map[string]any{"id": 8})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Delete, "/api/admin/remind/delete", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
	})
}
