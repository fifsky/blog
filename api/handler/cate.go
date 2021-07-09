package handler

import (
	"fmt"

	"app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/goapt/gee"
	"github.com/ilibs/gosql/v2"
)

type Cate struct {
	db       *gosql.DB
	cateRepo *repo.Cate
}

func NewCate(db *gosql.DB, cateRepo *repo.Cate) *Cate {
	return &Cate{db: db, cateRepo: cateRepo}
}

func (a *Cate) All(c *gee.Context) gee.Response {
	cates, err := a.cateRepo.GetAllCates()
	if err != nil {
		return response.Fail(c, 203, err)
	}

	data := make([]map[string]string, 0)

	for _, v := range cates {
		data = append(data, map[string]string{
			"url":     "/categroy/" + v.Domain,
			"content": fmt.Sprintf("%s(%d)", v.Name, v.Num),
		})
	}

	return response.Success(c, data)
}

func (a *Cate) List(c *gee.Context) gee.Response {
	cates, err := a.cateRepo.GetAllCates()
	if err != nil {
		return response.Fail(c, 203, err)
	}
	return response.Success(c, gee.H{
		"list":      cates,
		"pageTotal": len(cates),
	})
}

func (a *Cate) Post(c *gee.Context) gee.Response {
	cate := &model.Cates{}
	if err := c.ShouldBindJSON(cate); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if cate.Id > 0 {
		if _, err := a.db.Model(cate).Update(); err != nil {
			return response.Fail(c, 201, "更新失败")
		}
	} else {
		if _, err := a.db.Model(cate).Create(); err != nil {
			return response.Fail(c, 201, "创建失败")
		}
	}
	return response.Success(c, cate)
}

func (a *Cate) Delete(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	total, _ := a.db.Model(&model.Posts{}).Where("cate_id = ?", p.Id).Count()

	if total > 0 {
		return response.Fail(c, 201, "该分类下面还有文章，不能删除")
	}

	if _, err := a.db.Model(&model.Cates{Id: p.Id}).Delete(); err != nil {
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
