package handler

import (
	"app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/goapt/gee"
	"github.com/ilibs/gosql/v2"
)

type Link struct {
	db       *gosql.DB
	linkRepo *repo.Link
}

func NewLink(db *gosql.DB, linkRepo *repo.Link) *Link {
	return &Link{db: db, linkRepo: linkRepo}
}

func (l *Link) All(c *gee.Context) gee.Response {
	links, err := l.linkRepo.GetAllLinks()

	if err != nil {
		return response.Fail(c, 203, err)
	}

	data := make([]map[string]string, 0)

	for _, v := range links {
		data = append(data, map[string]string{
			"url":     v.Url,
			"content": v.Name,
		})
	}

	return response.Success(c, data)
}

func (l *Link) List(c *gee.Context) gee.Response {
	links, err := l.linkRepo.GetAllLinks()
	if err != nil {
		return response.Fail(c, 203, err)
	}

	return response.Success(c, gee.H{
		"list":      links,
		"pageTotal": len(links),
	})
}

func (l *Link) Post(c *gee.Context) gee.Response {
	link := &model.Links{}
	if err := c.ShouldBindJSON(link); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if link.Id > 0 {
		if _, err := l.db.Model(link).Update(); err != nil {
			return response.Fail(c, 201, "更新失败")
		}
	} else {
		if _, err := l.db.Model(link).Create(); err != nil {
			return response.Fail(c, 201, "创建失败")
		}
	}
	return response.Success(c, link)
}

func (l *Link) Delete(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if _, err := l.db.Model(&model.Links{Id: p.Id}).Delete(); err != nil {
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
