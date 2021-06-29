package handler

import (
	"app/provider/model"
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

	m := make(map[string]string)

	for k, v := range options.Map() {
		_, err := s.db.Model(&model.Options{
			OptionValue: v.String(),
		}).Where("option_key = ?", k).Update()

		if err != nil {
			return response.Fail(c, 203, err)
		}
		m[k] = v.String()
	}

	return response.Success(c, m)
}
