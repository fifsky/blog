package handler

import (
	"bytes"
	"net/http"
	"testing"

	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestLink_All(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(store.New(db))
		rr := doJSON(handler.All, "/api/link/all", map[string]any{})
		if rr.Code != http.StatusOK || rr.Body.Len() == 0 {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestLink_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(store.New(db))
		rr := doJSON(handler.List, "/api/admin/link/list", map[string]any{})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"list"`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestLink_Create(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(store.New(db))
		rr := doJSON(handler.Create, "/api/admin/link/create", map[string]any{"name": "demo", "url": "https://example.com", "desc": "demo"})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		// rr2 := doJSON(handler.Post, "/api/admin/link/post", map[string]any{"name": "demo", "desc": "demo"})
		// if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
		// 	t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		// }
	})
}

func TestLink_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("links")...)
		handler := NewLink(store.New(db))
		rr := doJSON(handler.Delete, "/api/admin/link/delete", map[string]any{"id": 3})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Delete, "/api/admin/link/delete", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
	})
}
