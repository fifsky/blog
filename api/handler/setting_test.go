package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/store"
	"app/testutil"

	"github.com/goapt/dbunit"
)

func TestSetting_Get(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		handler := NewSetting(store.New(db))
		req := httptest.NewRequest(http.MethodPost, "/api/setting", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Get(rr, req)
		if rr.Code != http.StatusOK || rr.Body.Len() == 0 {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}

func TestSetting_Post(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("options")...)
		handler := NewSetting(store.New(db))
		body, _ := json.Marshal(map[string]string{"key": "abc", "key2": "def"})
		req := httptest.NewRequest(http.MethodPost, "/api/admin/setting/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Post(rr, req)
		if rr.Code != http.StatusOK || !bytes.Contains(rr.Body.Bytes(), []byte(`"code":200`)) {
			t.Fatalf("unexpected: code=%d body=%s", rr.Code, rr.Body.String())
		}
	})
}
