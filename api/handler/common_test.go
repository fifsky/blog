package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app/model"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	type testCase struct {
		name     string
		body     string
		expected any
		wantErr  bool
	}

	cases := []testCase{
		{
			name:    "valid",
			body:    `{"user_name":"test","password":"123456"}`,
			wantErr: false,
			expected: LoginRequest{
				UserName: "test",
				Password: "123456",
			},
		},
		{
			name:    "empty",
			body:    `{"user_name":"","password":""}`,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := decode[LoginRequest](httptest.NewRequest(http.MethodPost, "/", strings.NewReader(c.body)))
			if c.wantErr {
				t.Log(err)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expected, v)
			}
		})
	}
}

func doJSON(handler func(http.ResponseWriter, *http.Request), url string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func doJSONWithUser(handler func(http.ResponseWriter, *http.Request), url string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	req = req.WithContext(context.WithValue(req.Context(), "userInfo", &model.User{Id: 1}))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}
