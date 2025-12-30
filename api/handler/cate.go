package handler

import (
	"fmt"
	"net/http"
	"time"

	"app/model"
	"app/response"
	"app/store"
)

type Cate struct {
	store *store.Store
}

func NewCate(s *store.Store) *Cate {
	return &Cate{store: s}
}

func (a *Cate) All(w http.ResponseWriter, r *http.Request) {
	cates, err := a.store.GetAllCates(r.Context())
	if err != nil {
		response.Fail(w, 203, err)
		return
	}

	data := make([]CateMenuItem, 0)

	for _, v := range cates {
		data = append(data, CateMenuItem{
			Url:     "/categroy/" + v.Domain,
			Content: fmt.Sprintf("%s(%d)", v.Name, v.Num),
		})
	}

	response.Success(w, data)
}

func (a *Cate) List(w http.ResponseWriter, r *http.Request) {
	cates, err := a.store.GetAllCates(r.Context())
	if err != nil {
		response.Fail(w, 203, err)
		return
	}
	items := make([]CateListItem, 0, len(cates))
	for _, c := range cates {
		items = append(items, CateListItem{
			Id:        c.Id,
			Name:      c.Name,
			Desc:      c.Desc,
			Domain:    c.Domain,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Num:       c.Num,
		})
	}
	resp := CateListResponse{
		List:      items,
		PageTotal: len(items),
	}
	response.Success(w, resp)
}

func (a *Cate) Post(w http.ResponseWriter, r *http.Request) {
	bodyCate, err := decode[CateRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	cate := bodyCate

	now := time.Now()
	if cate.Id > 0 {
		u := &model.UpdateCate{Id: cate.Id}
		if cate.Name != "" {
			u.Name = &cate.Name
		}
		if cate.Desc != "" {
			u.Desc = &cate.Desc
		}
		if cate.Domain != "" {
			u.Domain = &cate.Domain
		}
		u.UpdatedAt = &now
		if err := a.store.UpdateCate(r.Context(), u); err != nil {
			response.Fail(w, 201, "更新失败")
			return
		}
		response.Success(w, cate)
	}

	if cate.Name == "" || cate.Domain == "" {
		response.Fail(w, 201, "参数错误: 分类名或域名不能为空")
		return
	}
	c := &model.Cate{
		Name:      cate.Name,
		Desc:      cate.Desc,
		Domain:    cate.Domain,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := a.store.CreateCate(r.Context(), c); err != nil {
		response.Fail(w, 201, "创建失败")
		return
	}
	response.Success(w, cate)
}

func (a *Cate) Delete(w http.ResponseWriter, r *http.Request) {
	p, err := decode[IDRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	if p.Id <= 0 {
		response.Fail(w, 201, "参数错误: 分类ID不能为空")
		return
	}

	total, _ := a.store.PostsCount(r.Context(), p.Id)

	if total > 0 {
		response.Fail(w, 201, "该分类下面还有文章，不能删除")
		return
	}

	if err := a.store.DeleteCate(r.Context(), p.Id); err != nil {
		response.Fail(w, 201, "删除失败")
		return
	}
	response.Success(w, nil)
}
