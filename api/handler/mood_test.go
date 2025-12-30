package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/model"
	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestMood_List(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users")...)
		handler := NewMood(store.New(db))
		rr := doJSON(handler.List, "/api/admin/mood/list", map[string]any{"page": 1})
		if rr.Code != http.StatusOK || rr.Body.Len() == 0 {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestMood_Create(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("moods"))
		handler := NewMood(store.New(db))
		// success with user context
		b, _ := json.Marshal(map[string]any{"content": "demo"})
		req := httptest.NewRequest(http.MethodPost, "/api/admin/mood/create", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(SetLoginUser(req.Context(), &model.User{Id: 1}))
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}

		// rr2 := doJSON(handler.Post, "/api/admin/mood/post", map[string]any{})
		// if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
		// 	t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		// }
		//
		// rr3 := doJSON(handler.Post, "/api/admin/mood/post", map[string]any{"id": 1, "content": "demo2"})
		// if rr3.Code == http.StatusOK || !bytes.Contains(rr3.Body.Bytes(), []byte(`"code":200`)) {
		// 	t.Fatalf("unexpected: code=%d body=%s", rr3.Code, rr3.Body.String())
		// }
	})
}

func TestMood_Delete(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixture("moods"))
		handler := NewMood(store.New(db))
		rr := doJSON(handler.Delete, "/api/admin/mood/delete", map[string]any{"id": 1})
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
		rr2 := doJSON(handler.Delete, "/api/admin/mood/delete", map[string]any{})
		if rr2.Code == http.StatusOK || !bytes.Contains(rr2.Body.Bytes(), []byte(`"code":201`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr2.Code, rr2.Body.String())
		}
	})
}
