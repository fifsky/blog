package handler

import (
	"net/http"
	"time"

	"app/model"
	"app/response"
	"app/store"

	"github.com/samber/lo"
)

type Mood struct {
	store *store.Store
}

func NewMood(s *store.Store) *Mood {
	return &Mood{store: s}
}

func (m *Mood) List(w http.ResponseWriter, r *http.Request) {
	p, err := decode[PageRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	num := 10
	moods, err := m.store.ListMood(r.Context(), p.Page, num)
	if err != nil {
		response.Fail(w, 202, err)
		return
	}
	uids := lo.Map(moods, func(item model.Mood, index int) int {
		return item.UserId
	})

	um, err := m.store.GetUserByIds(r.Context(), uids)

	items := make([]MoodItem, 0, len(moods))
	for _, md := range moods {
		item := MoodItem{
			Id:        md.Id,
			Content:   md.Content,
			CreatedAt: md.CreatedAt,
		}
		if u, ok := um[md.UserId]; ok {
			item.User = &UserSummary{Id: u.Id, Name: u.Name, NickName: u.NickName}
		}
		items = append(items, item)
	}

	total, err := m.store.CountMoodTotal(r.Context())
	resp := MoodListResponse{
		List:      items,
		PageTotal: totalPages(total, num),
	}

	if err != nil {
		response.Fail(w, 500, err)
		return
	}

	response.Success(w, resp)
}

func (m *Mood) Post(w http.ResponseWriter, r *http.Request) {
	bodyMood, err := decode[MoodRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	in := bodyMood

	userId := 0
	if u := getLoginUser(r.Context()); u != nil {
		userId = u.Id
	}

	id := in.Id
	if id > 0 {
		u := &model.UpdateMood{Id: id}
		if in.Content != "" {
			u.Content = &in.Content
		}
		if err := m.store.UpdateMood(r.Context(), u); err != nil {
			response.Fail(w, 201, "更新心情失败")
			return
		}
	} else {
		c := &model.Mood{
			Content:   in.Content,
			UserId:    userId,
			CreatedAt: time.Now(),
		}
		if _, err := m.store.CreateMood(r.Context(), c); err != nil {
			response.Fail(w, 201, "发表心情失败")
			return
		}
	}

	response.Success(w, nil)
}

func (m *Mood) Delete(w http.ResponseWriter, r *http.Request) {
	p, err := decode[IDRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	if err := m.store.DeleteMood(r.Context(), p.Id); err != nil {
		response.Fail(w, 201, "删除失败")
		return
	}
	response.Success(w, nil)
}
