package contract

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// 用于测试解码的目标结构体
type decodeTarget struct {
	Page    int    `json:"page"`
	Keyword string `json:"keyword"`
	Debug   bool   `json:"debug"`
	Ids     []int  `json:"ids"`
	UserID  int    `json:"userID"`
	PostID  int    `json:"postID"`
}

func TestDecode_QueryParams(t *testing.T) {
	c := NewCodec()
	// 构造带有查询参数的 GET 请求
	req, err := http.NewRequest(http.MethodGet, "http://example.com?a=1&page=2&keyword=go&debug=true&ids=1&ids=2", nil)
	require.NoError(t, err)

	var v decodeTarget
	// 解码：应从 URL 查询参数填充到 v
	require.NoError(t, c.Decode(req, &v))

	require.Equal(t, 2, v.Page)
	require.Equal(t, "go", v.Keyword)
	require.Equal(t, true, v.Debug)
	require.Equal(t, []int{1, 2}, v.Ids)
}

func TestDecode_FormURLEncoded(t *testing.T) {
	c := NewCodec()
	// 构造 x-www-form-urlencoded 的 POST 请求
	body := "page=3&keyword=form&debug=false&ids=5&ids=6"
	req, err := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(body))
	require.NoError(t, err)
	// 设置表单类型头，便于 ParseForm 解析
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var v decodeTarget
	// 解码：应从表单参数填充到 v
	require.NoError(t, c.Decode(req, &v))

	require.Equal(t, 3, v.Page)
	require.Equal(t, "form", v.Keyword)
	require.Equal(t, false, v.Debug)
	require.Equal(t, []int{5, 6}, v.Ids)
}

func TestDecode_PathParams(t *testing.T) {
	c := NewCodec()
	// 构造带有路径模式与路径值的请求
	req, err := http.NewRequest(http.MethodGet, "http://example.com/users/42/posts/1001", nil)
	require.NoError(t, err)
	// 在请求上设置路径模式与具体变量值（Go 1.22+ 提供 PathValue/SetPathValue）
	req.Pattern = "/users/{userID}/posts/{postID}"
	req.SetPathValue("userID", "42")
	req.SetPathValue("postID", "1001")

	var v decodeTarget
	// 解码：应将路径变量注入到 v
	require.NoError(t, c.Decode(req, &v))

	require.Equal(t, 42, v.UserID)
	require.Equal(t, 1001, v.PostID)
}

func TestDecode_BodyJSON(t *testing.T) {
	c := NewCodec()
	// 构造 JSON 请求体
	payload := decodeTarget{
		Page:    7,
		Keyword: "json",
		Debug:   true,
		Ids:     []int{9, 10},
		UserID:  123,
	}
	buf, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(string(buf)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	var v decodeTarget
	// 解码：从 Body JSON 填充到 v
	require.NoError(t, c.Decode(req, &v))

	require.Equal(t, payload.Page, v.Page)
	require.Equal(t, payload.Keyword, v.Keyword)
	require.Equal(t, payload.Debug, v.Debug)
	require.Equal(t, payload.Ids, v.Ids)
	require.Equal(t, payload.UserID, v.UserID)
}
