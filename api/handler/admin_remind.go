package handler

import (
	"github.com/goapt/golib/pagination"

	"app/model"

	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

var AdminRemindList gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.Fail(201, "参数错误")
	}

	h := gin.H{}
	num := 10
	reminds, err := model.RemindGetList(p.Page, num)
	h["list"] = reminds

	total, err := gosql.Model(&model.Reminds{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return c.Fail(500, err)
	}

	return c.Success(h)
}

var AdminRemindPost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	remind := &model.Reminds{}
	if err := c.ShouldBindJSON(remind); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if remind.Content == "" {
		return c.Fail(201, "提醒内容不能为空")
	}

	if remind.Id > 0 {
		remind.Status = 1
		if _, err := gosql.Model(remind).Update(); err != nil {
			logger.Error(err)
			return c.Fail(201, "更新失败:"+err.Error())
		}
	} else {
		if _, err := gosql.Model(remind).Create(); err != nil {
			logger.Error(err)
			return c.Fail(201, "创建失败"+err.Error())
		}
	}
	return c.Success(remind)
}

var AdminRemindDelete gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.Fail(201, "参数错误")
	}

	if _, err := gosql.Model(&model.Reminds{Id: p.Id}).Delete(); err != nil {
		logger.Error(err)
		return c.Fail(201, "删除失败")
	}
	return c.Success(nil)
}
