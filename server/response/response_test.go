package response

import (
	"app/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	Success(w, map[string]any{"id": 1})
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusOK)
	}
	want := `{"id":1}`
	assert.JSONEq(t, want, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestFail(t *testing.T) {
	t.Run("msg", func(t *testing.T) {
		w := httptest.NewRecorder()
		Fail(w, errors.BadRequest("INVALID_PARAM", "参数错误"))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusBadRequest)
		}
		want := `{"code":"INVALID_PARAM","message":"参数错误"}`
		assert.JSONEq(t, want, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Fail(w, errors.InternalServer("SYSTEM_ERROR", "system error"))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusInternalServerError)
		}
		want := `{"code":"SYSTEM_ERROR","message":"system error"}`
		assert.JSONEq(t, want, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("other", func(t *testing.T) {
		w := httptest.NewRecorder()
		Fail(w, errors.InternalServer("SYSTEM_ERROR", "system error"))
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusInternalServerError)
		}
		want := `{"code":"SYSTEM_ERROR","message":"system error"}`
		assert.JSONEq(t, want, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}
