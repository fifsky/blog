package handler

import (
	"app/response"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/pagination"
	"github.com/ilibs/gosql/v2"

	"app/model"
)

var MoodList gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	h := gin.H{}
	num := 10
	moods, err := model.MoodGetList(p.Page, num)
	h["list"] = moods

	total, err := gosql.Model(&model.Moods{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, h)
}
