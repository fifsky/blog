package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// MockSuite is test set for http client
type MockSuite struct {
	URI                 string
	MatchBody           map[string]any
	MatchQuery          map[string]any
	ResponseBody        string
	StatusCode          int
	Header              http.Header
	Error               error
	ResponseInterceptor func(response *http.Response, body []byte) *http.Response
}

// Mock quickly define HTTP Response for mock RoundTriper
func Mock(suite []MockSuite) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			var g gjson.Result
			if strings.Contains(req.Header.Get("Content-Type"), "application/json") {
				var reqBody []byte
				if req.Body != nil {
					reqBody, _ = io.ReadAll(req.Body)
				}
				g = gjson.ParseBytes(reqBody)
			}

			if strings.Contains(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
				bb := make(map[string]any)
				_ = req.ParseForm()
				for cc, vv := range req.PostForm {
					if len(vv) == 1 {
						bb[cc] = vv[0]
					} else {
						bb[cc] = vv
					}
				}

				jb, _ := json.Marshal(bb)
				g = gjson.ParseBytes(jb)
			}

			for _, v := range suite {
				var re *regexp.Regexp
				if v.URI != "*" {
					// fmt.Println("===>",v.URI,req.URL.RequestURI())
					re = regexp.MustCompile(v.URI)
				}

				if v.URI == "*" || re.MatchString(req.URL.RequestURI()) {
					header := make(http.Header)
					if v.Header != nil {
						header = v.Header
					}

					if v.StatusCode == 0 {
						v.StatusCode = http.StatusOK
					}

					if v.MatchBody != nil {
						isMatchBody := true
						for rk, rv := range v.MatchBody {
							if g.Get(rk).String() != fmt.Sprint(rv) {
								isMatchBody = false
								break
							}
						}
						// 如果没有匹配到body的值，则跳过
						if !isMatchBody {
							continue
						}
					}

					if v.MatchQuery != nil {
						isMatchQuery := true
						query := req.URL.Query()
						for rk, rv := range v.MatchQuery {
							if query.Get(rk) != fmt.Sprint(rv) {
								isMatchQuery = false
								break
							}
						}
						// 如果没有匹配到query的值，则跳过
						if !isMatchQuery {
							continue
						}
					}

					if v.Error != nil {
						return nil, v.Error
					}
					return suitResponse(v, header)
				}
			}

			errResp := map[string]any{
				"error":        "HTTP Suite Miss",
				"request_uri":  req.URL.RequestURI(),
				"reqeust_body": g.String(),
			}

			errb, _ := json.Marshal(errResp)
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader(errb)),
				Header:     make(http.Header),
			}, nil
		})
	}
}

// 给返回值注入签名
func suitResponse(s MockSuite, header http.Header) (*http.Response, error) {
	if s.ResponseInterceptor != nil {
		return s.ResponseInterceptor(&http.Response{
			StatusCode: s.StatusCode,
			Header:     header,
		}, []byte(s.ResponseBody)), nil
	}

	return &http.Response{
		StatusCode: s.StatusCode,
		Body:       io.NopCloser(bytes.NewBufferString(s.ResponseBody)),
		Header:     header,
	}, nil
}
