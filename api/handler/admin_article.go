package handler

import (
	"net/http"
	"time"

	"app/response"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/hashing"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"

	"app/config"
	"app/model"
)

var AdminArticlePost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	post := &model.Posts{}
	if err := c.ShouldBindJSON(post); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	post.Status = 1
	post.UserId = getLoginUser(c).Id

	if post.Title == "" {
		return response.Fail(c, 201, "文章标题不能为空")
	}

	if post.CateId < 1 {
		return response.Fail(c, 201, "请选择文章分类")
	}

	if post.Id > 0 {
		if _, err := gosql.Model(post).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新文章失败")
		}
	} else {
		if _, err := gosql.Model(post).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "发表文章失败")
		}
	}

	return response.Success(c, post)
}

var AdminArticleDelete gee.HandlerFunc = func(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	post := &model.Posts{Id: p.Id, Status: 2}
	if _, err := gosql.Model(post).Update(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}

var AdminUploadPost gee.HandlerFunc = func(c *gee.Context) gee.Response {
	file, _, err := c.Request.FormFile("uploadFile")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.String("Bad request")
	}
	client, err := oss.New(config.App.OSS.Endpoint, config.App.OSS.AccessKey, config.App.OSS.AccessSecret)

	if err != nil {
		return c.JSON(gin.H{
			"errno": 201,
		})
	}

	bucket, _ := client.Bucket(config.App.OSS.Bucket)
	day := time.Now().Format("20060102")

	filename := "upload/" + day + "/" + hashing.Md5File(file) + ".png"
	file.Seek(0, 0)

	err = bucket.PutObject(filename, file)
	if err != nil {
		return c.JSON(gin.H{
			"errno": 202,
		})
	}

	return c.JSON(gin.H{
		"errno": 0,
		"data": [...]string{
			"https://static.fifsky.com/" + filename + "!blog",
		},
	})

}
