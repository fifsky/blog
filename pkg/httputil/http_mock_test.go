package httputil

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	t.Run("any", func(t *testing.T) {
		suites := []MockSuite{
			{
				URI:          "*",
				ResponseBody: "ok",
			},
		}

		client := NewClient(WithMiddleware(Mock(suites)))

		resp, err := client.Post("/post", "text/plain; charset=utf-8", nil)
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(body), "ok")
	})

	t.Run("custom", func(t *testing.T) {
		suites := []MockSuite{
			{
				URI:          "/get",
				ResponseBody: "ok1",
			},
			{
				URI:          "/user/id/.*",
				ResponseBody: "ok2",
			},
			{
				URI:          "/find\\?id=.*",
				ResponseBody: "ok3",
			},
			{
				URI:          "/bodymatch",
				ResponseBody: "ok4",
				MatchBody:    map[string]any{"user_id": 1},
			},
			{
				URI:   "/error",
				Error: errors.New("mock error"),
			},
			{
				URI:          "/query",
				ResponseBody: "ok5",
				MatchQuery:   map[string]any{"name": "test"},
			},
		}

		client := NewClient(WithMiddleware(Mock(suites)))

		uris := []struct {
			Uri          string
			StatusCode   int
			ResponseBody string
			HasError     bool
		}{
			{
				"/get",
				200,
				"ok1",
				false,
			}, {
				"/user/id/1",
				200,
				"ok2",
				false,
			}, {
				"/find?id=1",
				200,
				"ok3",
				false,
			}, {
				"/unkonw",
				404,
				`{"error":"HTTP Suite Miss","reqeust_body":"","request_uri":"/unkonw"}`,
				false,
			}, {
				"/error",
				404,
				"",
				true,
			}, {
				"/bodymatch",
				200,
				"ok4",
				false,
			}, {
				"/query?id=737373&name=test",
				200,
				"ok5",
				false,
			},
		}

		for _, v := range uris {
			t.Run(v.Uri, func(t *testing.T) {
				var body []byte
				ct := "text/plain; charset=utf-8"
				if v.Uri == "/bodymatch" {
					body = []byte(`{"user_id":1}`)
					ct = "application/json; charset=utf-8"
				}

				resp, err := client.Post(v.Uri, ct, bytes.NewReader(body))
				if v.HasError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				if resp != nil {
					assert.Equal(t, v.StatusCode, resp.StatusCode)
					if resp.Body != nil {
						body, err := io.ReadAll(resp.Body)
						assert.NoError(t, err)
						assert.Equal(t, v.ResponseBody, string(body))
					}
				}
			})
		}
	})
}

func Test_suitResponse(t *testing.T) {
	type args struct {
		s      MockSuite
		header http.Header
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test",
			args: args{
				s: MockSuite{
					StatusCode:   http.StatusOK,
					URI:          "/get",
					ResponseBody: "ok1",
				},
				header: make(http.Header),
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok1")),
				Header:     make(http.Header),
			},
		},
		{
			name: "test2",
			args: args{
				s: MockSuite{
					StatusCode: http.StatusOK,
					URI:        "/get",
					ResponseInterceptor: func(response *http.Response, body []byte) *http.Response {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok2")),
							Header:     make(http.Header),
						}
					},
				},
				header: make(http.Header),
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok2")),
				Header:     make(http.Header),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := suitResponse(tt.args.s, tt.args.header)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "suitResponse(%v, %v)", tt.args.s, tt.args.header)
		})
	}
}
