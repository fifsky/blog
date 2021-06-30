package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goapt/gee"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gee.CreateTestContext(w)
	resp := Success(ctx, gee.H{"id": 1})
	resp.Render()
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"code":200,"data":{"id":1},"msg":"success"}`, w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestFail(t *testing.T) {
	t.Run("msg", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gee.CreateTestContext(w)
		resp := Fail(ctx, 201, "参数错误")
		resp.Render()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"code":201,"msg":"参数错误"}`, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gee.CreateTestContext(w)
		resp := Fail(ctx, 500, errors.New("system error"))
		resp.Render()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"code":500,"msg":"system error"}`, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("other", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gee.CreateTestContext(w)
		resp := Fail(ctx, 500, gee.H{
			"error": "noterror",
		})
		resp.Render()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"code":500,"msg":"map[error:noterror]"}`, w.Body.String())
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}
