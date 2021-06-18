package handler

import (
	"fmt"

	"app/response"
	"github.com/goapt/gee"

	"app/model"
)

var CateAll gee.HandlerFunc = func(c *gee.Context) gee.Response {
	cates := model.GetAllCates()
	data := make([]map[string]string, 0)

	for _, v := range cates {
		data = append(data, map[string]string{
			"url":     "/categroy/" + v.Domain,
			"content": fmt.Sprintf("%s(%d)", v.Name, v.Num),
		})
	}

	return response.Success(c, data)
}
