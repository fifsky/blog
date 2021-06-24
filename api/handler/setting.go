package handler

import (
	model2 "app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/goapt/gee"
	"github.com/ilibs/gosql/v2"
	"github.com/tidwall/gjson"
)

type Setting struct {
	db          *gosql.DB
	settingRepo *repo.Setting
}

func NewSetting(db *gosql.DB, settingRepo *repo.Setting) *Setting {
	return &Setting{db: db, settingRepo: settingRepo}
}

func (s *Setting) Get(c *gee.Context) gee.Response {
	m, err := s.settingRepo.GetOptions()
	if err != nil {
		return response.Fail(c, 202, err)
	}

	return response.Success(c, m)
}

func (s *Setting) Post(c *gee.Context) gee.Response {

	body, err := c.GetRawData()
	if err != nil {
		return response.Fail(c, 202, err)
	}

	options := gjson.ParseBytes(body)

	for k, v := range options.Map() {
		s.db.Model(&model2.Options{
			OptionValue: v.String(),
		}).Where("option_key = ?", k).Update()
	}

	return response.Success(c, options)
}
