package handler

import (
	"time"

	model2 "app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/pagination"
	"github.com/goapt/golib/robot"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type Remind struct {
	db         *gosql.DB
	remindRepo *repo.Remind
}

func NewRemind(db *gosql.DB, remindRepo *repo.Remind) *Remind {
	return &Remind{db: db, remindRepo: remindRepo}
}

func (r *Remind) Change(c *gee.Context) gee.Response {
	remind, ok := c.Get("remind")

	if !ok {
		return response.Fail(c, 202, "记录未找到")
	}

	v := remind.(*model2.Reminds)

	_, err := r.db.Model(&model2.Reminds{Status: 1}).Where("id = ?", v.Id).Update()

	if err != nil {
		return response.Fail(c, 203, err)
	}

	_ = robot.Message("已确认收到提醒")
	return c.String("已确认收到提醒")
}

func (r *Remind) Delay(c *gee.Context) gee.Response {
	remind, ok := c.Get("remind")

	if !ok {
		return response.Fail(c, 202, "记录未找到")
	}

	v := remind.(*model2.Reminds)

	nextTime := time.Now().Add(10 * time.Minute)
	_, err := r.db.Model(&model2.Reminds{NextTime: nextTime}).Where("id = ?", v.Id).Update()

	if err != nil {
		return response.Fail(c, 203, err)
	}
	_ = robot.Message("将在10分钟后再次提醒")
	return c.String("将在10分钟后再次提醒")
}

func (r *Remind) List(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误")
	}

	h := gin.H{}
	num := 10
	reminds, err := r.remindRepo.RemindGetList(p.Page, num)
	h["list"] = reminds

	total, err := r.db.Model(&model2.Reminds{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, h)
}

func (r *Remind) Post(c *gee.Context) gee.Response {
	remind := &model2.Reminds{}
	if err := c.ShouldBindJSON(remind); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if remind.Content == "" {
		return response.Fail(c, 201, "提醒内容不能为空")
	}

	if remind.Id > 0 {
		remind.Status = 1
		if _, err := r.db.Model(remind).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新失败:"+err.Error())
		}
	} else {
		if _, err := r.db.Model(remind).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "创建失败"+err.Error())
		}
	}
	return response.Success(c, remind)
}

func (r *Remind) Delete(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误")
	}

	if _, err := r.db.Model(&model2.Reminds{Id: p.Id}).Delete(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
