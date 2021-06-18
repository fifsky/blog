package handler

import (
	"app/response"
	"github.com/goapt/gee"
	"github.com/ilibs/gosql/v2"
	"github.com/tidwall/gjson"

	"app/model"
)

var AdminSetting gee.HandlerFunc = func(c *gee.Context) gee.Response {
	m, err := model.GetOptions()
	if err != nil {
		response.Fail(c, 202, err)
	}

	return response.Success(c, m)
}

var AdminSettingPost gee.HandlerFunc = func(c *gee.Context) gee.Response {

	body, err := c.GetRawData()
	if err != nil {
		return response.Fail(c, 202, err)
	}

	options := gjson.ParseBytes(body)

	for k, v := range options.Map() {
		gosql.Model(&model.Options{
			OptionValue: v.String(),
		}).Where("option_key = ?", k).Update()
	}

	return response.Success(c, options)
}
