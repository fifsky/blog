package handler

import (
	"bytes"
	"net/http"
	"testing"

	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestCate_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		handler := NewCate(store.New(db))
		rr := doJSON(handler.All, "/api/cate/all", map[string]any{})
		if rr.Code != http.StatusOK || rr.Body.Len() == 0 {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestCate_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		handler := NewCate(store.New(db))
		rr := doJSON(handler.List, "/api/admin/cate/list", map[string]any{})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"list"`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestCate_CreateUpdate(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates")...)
		handler := NewCate(store.New(db))

		rr := doJSON(handler.Create, "/api/admin/cate/create", map[string]any{"name": "demo", "domain": "demo", "desc": "demo"})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Update, "/api/admin/cate/update", map[string]any{"id": 1, "domain": "test2", "name": "test2", "desc": "test2", "updated_at": "2021-06-29 11:55:09"})
		if rr2.Code != http.StatusOK {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
		rr3 := doJSON(handler.Create, "/api/admin/cate/create", map[string]any{"name": "demo", "desc": "demo"})
		if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		}
	})
}

func TestCate_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates", "posts")...)
		handler := NewCate(store.New(db))
		rr := doJSON(handler.Delete, "/api/admin/cate/delete", map[string]any{"id": 3})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Delete, "/api/admin/cate/delete", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
		rr3 := doJSON(handler.Delete, "/api/admin/cate/delete", map[string]any{"id": 1})
		if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`不能删除`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		}
	})
}
