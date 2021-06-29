package handler

import (
	"app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/pagination"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type Mood struct {
	db       *gosql.DB
	moodRepo *repo.Mood
}

func NewMood(db *gosql.DB, moodRepo *repo.Mood) *Mood {
	return &Mood{db: db, moodRepo: moodRepo}
}

func (m *Mood) List(c *gee.Context) gee.Response {
	p := &struct {
		Page int `json:"page" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	h := gin.H{}
	num := 10
	moods, err := m.moodRepo.MoodGetList(p.Page, num)
	h["list"] = moods

	total, err := m.db.Model(&model.Moods{}).Count()
	pager := pagination.New(int(total), num, p.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}

	return response.Success(c, h)
}

func (m *Mood) Post(c *gee.Context) gee.Response {
	mood := &model.Moods{}
	if err := c.ShouldBindJSON(mood); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	mood.UserId = getLoginUser(c).Id

	if mood.Id > 0 {
		if _, err := m.db.Model(mood).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新心情失败")
		}
	} else {
		if _, err := m.db.Model(mood).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "发表心情失败")
		}
	}

	return response.Success(c, nil)
}

func (m *Mood) Delete(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	if _, err := m.db.Model(&model.Moods{Id: p.Id}).Delete(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}
