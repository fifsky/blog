package testutil

import (
	"net/http"
	"testing"

	"github.com/goapt/gee"
	"github.com/goapt/test"
	"github.com/stretchr/testify/require"
)

type AssertType int

const (
	AssertEqual AssertType = iota
	AssertContains
)

type TestCase struct {
	Name         string
	RequestBody  interface{}
	ResponseBody string
	AssertType   AssertType
	Check        func(t *testing.T, response *test.Response)
	WantCode     int
}

func (tc *TestCase) Run(t *testing.T, url string, handlers ...gee.HandlerFunc) {
	t.Run(tc.Name, func(t *testing.T) {
		req := test.NewRequest(url, handlers...)
		resp, err := req.JSON(tc.RequestBody)
		require.NoError(t, err)
		if tc.WantCode > 0 {
			require.Equal(t, tc.WantCode, resp.Code)
		} else {
			require.Equal(t, http.StatusOK, resp.Code)
		}

		if tc.Check != nil {
			tc.Check(t, resp)
			return
		}

		switch tc.AssertType {
		case AssertEqual:
			require.Equal(t, tc.ResponseBody, resp.GetBodyString())
		case AssertContains:
			require.Contains(t, resp.GetBodyString(), tc.ResponseBody)
		}
	})
}
