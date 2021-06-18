package handler

import (
	"fmt"

	"app/response"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/convert"
	"github.com/goapt/golib/pagination"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"

	"app/model"
)

var ArchiveArticle gee.HandlerFunc = func(c *gee.Context) gee.Response {
	archives, err := model.PostArchive()
	if err != nil {
		logger.Error(err)
	}
	data := make([]map[string]string, 0)

	for _, v := range archives {
		data = append(data, map[string]string{
			"url":     "/date/" + v["ym"],
			"content": fmt.Sprintf("%s(%s)", v["ym"], v["total"]),
		})
	}

	return response.Success(c, data)
}

var ListArticle gee.HandlerFunc = func(c *gee.Context) gee.Response {
	options, err := model.GetOptions()
	if err != nil {
		response.Fail(c, 202, err)
	}

	num, err := convert.StrTo(options["post_num"]).Int()
	if err != nil || num < 1 {
		num = 10
	}

	req := &struct {
		Year    string `json:"year"`
		Month   string `json:"month"`
		Domain  string `json:"domain"`
		Keyword string `json:"keyword"`
		Page    int    `json:"page"`
		Type    int    `json:"type"`
	}{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, 201, "参数错误")
	}

	cate := &model.Cates{}

	if req.Domain != "" {
		cate.Domain = req.Domain
		gosql.Model(cate).Get()
	}

	artdate := ""

	if req.Year != "" && req.Month != "" {
		artdate = req.Year + "-" + req.Month
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	post := &model.Posts{}
	if cate.Id > 0 {
		post.CateId = cate.Id
	}

	posts, err := model.PostGetList(post, page, num, artdate, req.Keyword)
	if err != nil {
		return response.Fail(c, 500, err)
	}

	h := gin.H{}
	h["list"] = posts

	builder := gosql.Model(post)

	if artdate != "" {
		builder.Where("DATE_FORMAT(created_at,'%Y-%m') = ?", artdate)
	}

	if req.Keyword != "" {
		builder.Where("title like ?", "%"+req.Keyword+"%")
	}

	total, err := builder.Count()
	pager := pagination.New(int(total), num, req.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}
	return response.Success(c, h)
}

var PrevNextArticle gee.HandlerFunc = func(c *gee.Context) gee.Response {
	req := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(req); err != nil {
		return response.Fail(c, 201, "参数错误")
	}

	h := gin.H{}
	h["prev"] = map[string]interface{}{}
	h["next"] = map[string]interface{}{}

	prev, err := model.PostPrev(req.Id)
	if err == nil {
		h["prev"] = gin.H{
			"id":    prev.Id,
			"title": prev.Title,
		}
	}
	next, err := model.PostNext(req.Id)
	if err == nil {
		h["next"] = gin.H{
			"id":    next.Id,
			"title": next.Title,
		}
	}
	return response.Success(c, h)
}

var DetailArticle gee.HandlerFunc = func(c *gee.Context) gee.Response {
	req := &struct {
		Id  int    `json:"id"`
		Url string `json:"url"`
	}{}
	if err := c.ShouldBindJSON(req); err != nil {
		return response.Fail(c, 201, "参数错误")
	}

	post := &model.UserPosts{}

	if req.Id > 0 {
		post.Id = req.Id
	}

	if req.Url != "" {
		post.Url = req.Url
	}

	err := gosql.Model(post).Where("status = 1").Get()
	if err != nil {
		return response.Fail(c, 202, "您访问的文章不存在或已经删除！")
	}

	return response.Success(c, post)
}
