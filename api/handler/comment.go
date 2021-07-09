package handler

import (
	"fmt"

	"app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/golib/pagination"
	"github.com/goapt/golib/robot"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type Comment struct {
	db          *gosql.DB
	commentRepo *repo.Comment
}

func NewComment(db *gosql.DB, commentRepo *repo.Comment) *Comment {
	return &Comment{db: db, commentRepo: commentRepo}
}

func (m *Comment) List(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	comments, err := m.commentRepo.PostComments(p.Id, 1, 100)
	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, comments)
}

func (m *Comment) Post(c *gee.Context) gee.Response {
	comment := &model.Comments{}
	if err := c.ShouldBindJSON(comment); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	post := &model.Posts{}
	err := m.db.Model(post).Where("id = ?", comment.PostId).Get()
	if err != nil {
		return response.Fail(c, 201, "文章不存在")
	}

	comment.IP = c.ClientIP()

	if _, err := m.db.Model(comment).Create(); err != nil {
		return response.Fail(c, 201, "评论失败"+err.Error())
	}

	content := "您有新的评论!\n"
	content += fmt.Sprintf("文章:%s\n", post.Title)
	content += fmt.Sprintf("评论内容:%s\n", comment.Content)
	content += fmt.Sprintf("评论昵称:%s\n", comment.Name)
	content += fmt.Sprintf("评论时间:%s\n", comment.CreatedAt.Format("2006-01-02 15:04:05"))
	content += fmt.Sprintf("评论IP:%s\n", comment.IP)

	_ = robot.Message(content)
	return response.Success(c, comment)
}

func (m *Comment) Top(c *gee.Context) gee.Response {
	comments, err := m.commentRepo.NewComments()
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

	return response.Success(c, data)
}

func (m *Comment) AdminList(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	h := gee.H{}
	num := 10
	comments, err := m.commentRepo.CommentList(p.Page, num)
	h["list"] = comments

	total, err := m.db.Model(&model.Comments{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, h)
}

func (m *Comment) Delete(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if _, err := m.db.Model(&model.Comments{Id: p.Id}).Delete(); err != nil {
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
