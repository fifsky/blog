package handler

import (
	"app/model"
	"app/response"

	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

var AdminMoodPost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	mood := &model.Moods{}
	if err := c.ShouldBindJSON(mood); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	mood.UserId = getLoginUser(c).Id

	if mood.Content == "" {
		return response.Fail(c, 201, "内容不能为空")
	}

	if mood.Id > 0 {
		if _, err := gosql.Model(mood).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新心情失败")
		}
	} else {
		if _, err := gosql.Model(mood).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "发表心情失败")
		}
	}

	return response.Success(c, nil)
}

var AdminMoodDelete gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	if _, err := gosql.Model(&model.Moods{Id: p.Id}).Delete(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
