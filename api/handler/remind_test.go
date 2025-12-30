package handler

import (
	"bytes"
	"net/http"
	"testing"

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
		rr := doJSONWithRemind(handler.Change, "/api/remind/change", map[string]any{"id": 8})
		if rr.Code != http.StatusOK || rr.Body.String() == "" {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestRemind_Delay(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		robot := wechat.NewRobot("123")
		handler := NewRemind(store.New(db), robot)
		rr := doJSONWithRemind(handler.Delay, "/api/remind/delay", map[string]any{"id": 8})
		if rr.Code != http.StatusOK || rr.Body.String() == "" {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
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

func TestRemind_Create(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds")...)
		handler := NewRemind(store.New(db), wechat.NewRobot("123"))
		rr := doJSON(handler.Create, "/api/admin/remind/create", map[string]any{"type": 1, "content": "demo", "month": 1, "week": 0, "day": 1, "hour": 1, "minute": 1, "status": 0, "created_at": "2021-06-29 11:55:09"})
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
