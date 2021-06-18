package handler

import (
	"app/response"
	"github.com/goapt/gee"

	"app/model"
)

var LinkAll gee.HandlerFunc = func(c *gee.Context) gee.Response {
	links := model.GetAllLinks()
	data := make([]map[string]string, 0)

	for _, v := range links {
		data = append(data, map[string]string{
			"url":     v.Url,
			"content": v.Name,
		})
	}

	return response.Success(c, data)
}
