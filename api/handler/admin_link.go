package handler

import (
	"app/model"
	"app/response"

	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

var AdminLinkList gee.HandlerFunc = func(c *gee.Context) gee.Response {
	links := model.GetAllLinks()
	return response.Success(c, gin.H{
		"list":      links,
		"pageTotal": len(links),
	})
}

var AdminLinkPost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	link := &model.Links{}
	if err := c.ShouldBindJSON(link); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if link.Name == "" {
		return response.Fail(c, 201, "连接名称不能为空")
	}

	if link.Url == "" {
		return response.Fail(c, 201, "连接地址不能为空")
	}

	if link.Id > 0 {
		if _, err := gosql.Model(link).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(link).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "创建失败")
		}
	}
	return response.Success(c, link)
}

var AdminLinkDelete gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	if _, err := gosql.Model(&model.Links{Id: p.Id}).Delete(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
