package handler

import (
	"fmt"
	"net/http"
	"time"

	"app/model"
	"app/response"
	"app/store"
)

type Link struct {
	store *store.Store
}

func NewLink(s *store.Store) *Link {
	return &Link{store: s}
}

func (l *Link) All(w http.ResponseWriter, r *http.Request) {
	links, err := l.store.GetAllLinks(r.Context())

	if err != nil {
		response.Fail(w, 203, err)
		return
	}

	data := make([]LinkMenuItem, 0)

	for _, v := range links {
		data = append(data, LinkMenuItem{
			Url:     v.Url,
			Content: v.Name,
		})
	}

	response.Success(w, data)
}

func (l *Link) List(w http.ResponseWriter, r *http.Request) {
	links, err := l.store.GetAllLinks(r.Context())
	if err != nil {
		response.Fail(w, 203, err)
		return
	}

	items := make([]LinkItem, 0, len(links))
	for _, v := range links {
		items = append(items, LinkItem{
			Id:        v.Id,
			Name:      v.Name,
			Url:       v.Url,
			Desc:      v.Desc,
			CreatedAt: v.CreatedAt,
		})
	}
	resp := LinkListResponse{
		List:      items,
		PageTotal: len(items),
	}
	response.Success(w, resp)
}

func (l *Link) Post(w http.ResponseWriter, r *http.Request) {
	bodyLink, err := decode[LinkRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	in := bodyLink

	id := in.Id
	if id > 0 {
		u := &model.UpdateLink{Id: id, Name: nil, Url: nil, Desc: nil}
		if in.Name != "" {
			u.Name = &in.Name
		}
		if in.Url != "" {
			u.Url = &in.Url
		}
		if in.Desc != "" {
			u.Desc = &in.Desc
		}
		if err := l.store.UpdateLink(r.Context(), u); err != nil {
			response.Fail(w, 201, "更新失败")
			return
		}
	} else {
		c := &model.Link{
			Name:      in.Name,
			Url:       in.Url,
			Desc:      in.Desc,
			CreatedAt: time.Now(),
		}
		if _, err := l.store.CreateLink(r.Context(), c); err != nil {
			response.Fail(w, 201, fmt.Sprintf("创建失败: %v", err))
			return
		}
	}
	response.Success(w, in)
}

func (l *Link) Delete(w http.ResponseWriter, r *http.Request) {
	p, err := decode[IDRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	if err := l.store.DeleteLink(r.Context(), p.Id); err != nil {
		response.Fail(w, 201, "删除失败")
		return
	}
	response.Success(w, nil)
}
