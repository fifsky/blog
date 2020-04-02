package handler

import (
	"fmt"
	"time"

	"github.com/goapt/gee"
	"github.com/goapt/golib/robot"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"

	"app/model"
)

var CommentList gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.Fail(201, "参数错误")
	}

	comments, err := model.PostComments(p.Id, 1, 100)
	if err != nil {
		c.Fail(500, err)
	}

	return c.Success(comments)
}

var CommentPost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	return c.Fail(201, "该功能已下线，请按 Ctrl + F5 强制刷新页面！")

	comment := &model.Comments{}
	if err := c.ShouldBindJSON(comment); err != nil {
		return c.Fail(201, "参数错误:"+err.Error())
	}

	if comment.Name == "" {
		return c.Fail(201, "昵称不能为空")
	}

	if comment.Content == "" {
		return c.Fail(201, "评论内容不能为空")
	}

	if comment.PostId <= 0 {
		return c.Fail(201, "非法评论")
	}

	//if err := TCaptchaVerify(c.PostForm("ticket"), c.PostForm("randstr"), c.ClientIP()); err != nil {
	//	return c.Fail(201, err)
	//}

	post := &model.Posts{}
	err := gosql.Model(post).Where("id = ?", comment.PostId).Get()
	if err != nil {
		return c.Fail(201, "文章不存在")
	}

	comment.CreatedAt = time.Now()
	comment.IP = c.ClientIP()

	if _, err := gosql.Model(comment).Create(); err != nil {
		logger.Error(err)
		return c.Fail(201, "评论失败"+err.Error())
	}

	content := "您有新的评论!\n"
	content += fmt.Sprintf("文章:%s\n", post.Title)
	content += fmt.Sprintf("评论内容:%s\n", comment.Content)
	content += fmt.Sprintf("评论昵称:%s\n", comment.Name)
	content += fmt.Sprintf("评论时间:%s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("评论IP:%s\n", comment.IP)

	_ = robot.Message(content)
	return c.Success(comment)
}

var NewComment gee.HandlerFunc = func(c *gee.Context) gee.Response {
	comments, err := model.NewComments()
	if err != nil {
		logger.Error(err)
	}
	data := make([]map[string]string, 0)

	for _, v := range comments {
		var url string
		if v.Type == 2 {
			url = "/" + v.Url + "#comments"
		} else {
			url = fmt.Sprintf("/article/%d#comments", v.PostId)
		}

		data = append(data, map[string]string{
			"url":     url,
			"content": v.Content,
		})
	}

	return c.Success(data)
}
