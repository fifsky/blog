package handler

import (
	"net/http"
	"time"

	"app/model"
	"app/pkg/wechat"
	"app/response"
	"app/store"
)

type Remind struct {
	store *store.Store
	robot *wechat.Robot
}

func NewRemind(s *store.Store, robot *wechat.Robot) *Remind {
	return &Remind{
		store: s,
		robot: robot,
	}
}

func (r *Remind) Change(w http.ResponseWriter, req *http.Request) {
	remind := getRemind(req.Context())
	err := r.store.UpdateRemindStatus(req.Context(), remind.Id, 1)
	if err != nil {
		response.Fail(w, 203, err)
		return
	}

	_ = r.robot.Message("已确认收到提醒")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("已确认收到提醒"))
}

func (r *Remind) Delay(w http.ResponseWriter, req *http.Request) {
	remind := getRemind(req.Context())
	if remind == nil {
		response.Fail(w, 202, "记录未找到")
		return
	}
	err := r.store.UpdateRemindNextTime(req.Context(), remind.Id, remind.NextTime)

	if err != nil {
		response.Fail(w, 203, err)
		return
	}
	_ = r.robot.Message("将在10分钟后再次提醒")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("将在10分钟后再次提醒"))
}

func (r *Remind) List(w http.ResponseWriter, req *http.Request) {
	p, err := decode[PageRequest](req)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	num := 10
	reminds, err := r.store.ListRemind(req.Context(), p.Page, num)
	if err != nil {
		response.Fail(w, 202, err)
		return
	}
	items := make([]RemindItem, 0, len(reminds))
	for _, v := range reminds {
		items = append(items, RemindItem{
			Id:        v.Id,
			Type:      v.Type,
			Content:   v.Content,
			Month:     v.Month,
			Week:      v.Week,
			Day:       v.Day,
			Hour:      v.Hour,
			Minute:    v.Minute,
			Status:    v.Status,
			NextTime:  v.NextTime,
			CreatedAt: v.CreatedAt,
		})
	}

	total, err := r.store.CountRemindTotal(req.Context())
	resp := RemindListResponse{
		List:      items,
		PageTotal: totalPages(total, num),
	}

	if err != nil {
		response.Fail(w, 500, err)
		return
	}

	response.Success(w, resp)
}

func (r *Remind) Create(w http.ResponseWriter, req *http.Request) {
	in, err := decode[RemindCreateRequest](req)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	c := &model.Remind{
		Type:      in.Type,
		Content:   in.Content,
		Month:     in.Month,
		Week:      in.Week,
		Day:       in.Day,
		Hour:      in.Hour,
		Minute:    in.Minute,
		Status:    1,
		CreatedAt: time.Now(),
	}
	var lastId int64
	if lastId, err = r.store.CreateRemind(req.Context(), c); err != nil {
		response.Fail(w, 201, "创建失败"+err.Error())
		return
	}
	response.Success(w, IDResponse{Id: int(lastId)})
}

func (r *Remind) Update(w http.ResponseWriter, req *http.Request) {
	in, err := decode[RemindUpdateRequest](req)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}
	u := &model.UpdateRemind{Id: in.Id}
	if in.Type > 0 {
		u.Type = &in.Type
	}
	if in.Content != "" {
		u.Content = &in.Content
	}
	if in.Month > 0 {
		u.Month = &in.Month
	}
	if in.Week > 0 {
		u.Week = &in.Week
	}
	if in.Day > 0 {
		u.Day = &in.Day
	}
	if in.Hour > 0 {
		u.Hour = &in.Hour
	}
	if in.Minute > 0 {
		u.Minute = &in.Minute
	}
	if err := r.store.UpdateRemind(req.Context(), u); err != nil {
		response.Fail(w, 201, "更新失败:"+err.Error())
		return
	}
	response.Success(w, IDResponse{Id: in.Id})
}

func (r *Remind) Delete(w http.ResponseWriter, req *http.Request) {
	p, err := decode[IDRequest](req)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	if err := r.store.DeleteRemind(req.Context(), p.Id); err != nil {
		response.Fail(w, 201, "删除失败")
		return
	}
	response.Success(w, nil)
}
