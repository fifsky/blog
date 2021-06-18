package response

import (
	"fmt"

	"github.com/goapt/gee"
)

type ApiResponse struct {
	Context *gee.Context `json:"-"`
	Code    int          `json:"code"`
	Data    interface{}  `json:"data"`
	Msg     string       `json:"msg"`
}

func (c *ApiResponse) Render() {
	c.Context.JSON(c).Render()
}

func Success(c *gee.Context, data interface{}) gee.Response {

	return &ApiResponse{
		Context: c,
		Code:    200,
		Data:    data,
		Msg:     "success",
	}
}

func Fail(c *gee.Context, code int, msg interface{}) gee.Response {
	var m string
	switch e := msg.(type) {
	case string:
		m = e
	case error:
		m = e.Error()
	default:
		m = fmt.Sprintf("[%d]%v", code, e)
	}

	return &ApiResponse{
		Context: c,
		Code:    code,
		Data:    nil,
		Msg:     m,
	}
}
