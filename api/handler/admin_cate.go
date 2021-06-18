package handler

import (
	"app/model"
	"app/response"

	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

var AdminCateList gee.HandlerFunc = func(c *gee.Context) gee.Response {
	cates := model.GetAllCates()
	return response.Success(c, gin.H{
		"list":      cates,
		"pageTotal": len(cates),
	})
}

var AdminCatePost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	cate := &model.Cates{}
	if err := c.ShouldBindJSON(cate); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if cate.Name == "" {
		return response.Fail(c, 201, "分类名不能为空")
	}

	if cate.Domain == "" {
		return response.Fail(c, 201, "分类缩略名不能为空")
	}

	if cate.Id > 0 {
		if _, err := gosql.Model(cate).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新失败")
		}
	} else {
		if _, err := gosql.Model(cate).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "创建失败")
		}
	}
	return response.Success(c, cate)
}

var AdminCateDelete gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	total, _ := gosql.Model(&model.Posts{}).Where("cate_id = ?", p.Id).Count()

	if total > 0 {
		return response.Fail(c, 201, "该分类下面还有文章，不能删除")
	}

	if _, err := gosql.Model(&model.Cates{Id: p.Id}).Delete(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
