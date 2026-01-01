package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	Success(w, map[string]any{"id": 1})
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusOK)
	}
	want := `{"id":1}` + "\n"
	if w.Body.String() != want {
		t.Fatalf("unexpected body: got=%s want=%s", w.Body.String(), want)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content-type: got=%s want=%s", ct, "application/json; charset=utf-8")
	}
}

func TestFail(t *testing.T) {
	t.Run("msg", func(t *testing.T) {
		w := httptest.NewRecorder()
		Fail(w, 201, "参数错误")
		if w.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusBadRequest)
		}
		want := `{"code":201,"msg":"参数错误"}`
		if w.Body.String() != want {
			t.Fatalf("unexpected body: got=%s want=%s", w.Body.String(), want)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Fatalf("unexpected content-type: got=%s want=%s", ct, "application/json; charset=utf-8")
		}
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Fail(w, 500, errors.New("system error"))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusBadRequest)
		}
		want := `{"code":500,"msg":"system error"}`
		if w.Body.String() != want {
			t.Fatalf("unexpected body: got=%s want=%s", w.Body.String(), want)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Fatalf("unexpected content-type: got=%s want=%s", ct, "application/json; charset=utf-8")
		}
	})

	t.Run("other", func(t *testing.T) {
		w := httptest.NewRecorder()
		Fail(w, 500, map[string]string{"error": "noterror"})
		if w.Code != http.StatusBadRequest {
			t.Fatalf("unexpected status code: got=%d want=%d", w.Code, http.StatusBadRequest)
		}
		want := `{"code":500,"msg":"map[error:noterror]"}`
		if w.Body.String() != want {
			t.Fatalf("unexpected body: got=%s want=%s", w.Body.String(), want)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Fatalf("unexpected content-type: got=%s want=%s", ct, "application/json; charset=utf-8")
		}
	})
}
