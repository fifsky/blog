package handler

import (
	"github.com/goapt/gee"
	"github.com/ilibs/gosql/v2"
	"github.com/tidwall/gjson"

	"app/model"
)

var AdminSetting gee.HandlerFunc = func(c *gee.Context) gee.Response {
	m, err := model.GetOptions()
	if err != nil {
		c.Fail(202, err)
	}

	return c.Success(m)
}

var AdminSettingPost gee.HandlerFunc = func(c *gee.Context) gee.Response {

	body, err := c.GetRawData()
	if err != nil {
		return c.Fail(202, err)
	}

	options := gjson.ParseBytes(body)

	for k, v := range options.Map() {
		gosql.Model(&model.Options{
			OptionValue: v.String(),
		}).Where("option_key = ?", k).Update()
	}

	return c.Success(options)
}
