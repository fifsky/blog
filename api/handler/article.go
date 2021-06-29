package handler

import (
	"fmt"
	"net/http"
	"time"

	"app/config"
	"app/provider/model"
	"app/provider/repo"
	"app/response"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/goapt/gee"
	"github.com/goapt/golib/convert"
	"github.com/goapt/golib/hashing"
	"github.com/goapt/golib/pagination"
	"github.com/goapt/logger"
	"github.com/gorilla/feeds"
	"github.com/ilibs/gosql/v2"
)

type Article struct {
	db          *gosql.DB
	artRepo     *repo.Article
	settingRepo *repo.Setting
}

func NewArticle(db *gosql.DB, artRepo *repo.Article, settingRepo *repo.Setting) *Article {
	return &Article{db: db, artRepo: artRepo, settingRepo: settingRepo}
}

func (a *Article) Archive(c *gee.Context) gee.Response {
	archives, err := a.artRepo.PostArchive()
	if err != nil {
		logger.Error(err)
	}
	data := make([]map[string]string, 0)

	for _, v := range archives {
		data = append(data, map[string]string{
			"url":     "/date/" + v["ym"],
			"content": fmt.Sprintf("%s(%s)", v["ym"], v["total"]),
		})
	}

	return response.Success(c, data)
}

func (a *Article) List(c *gee.Context) gee.Response {
	options, err := a.settingRepo.GetOptions()
	if err != nil {
		return response.Fail(c, 202, err)
	}

	num, err := convert.StrTo(options["post_num"]).Int()
	if err != nil || num < 1 {
		num = 10
	}

	req := &struct {
		Year    string `json:"year"`
		Month   string `json:"month"`
		Domain  string `json:"domain"`
		Keyword string `json:"keyword"`
		Page    int    `json:"page" binding:"required"`
		Type    int    `json:"type"`
	}{}
	if err := c.ShouldBindJSON(req); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	cate := &model.Cates{}

	if req.Domain != "" {
		cate.Domain = req.Domain
		err := a.db.Model(cate).Get()
		if err != nil {
			return response.Fail(c, 202, err)
		}
	}

	artdate := ""

	if req.Year != "" && req.Month != "" {
		artdate = req.Year + "-" + req.Month
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	post := &model.Posts{}
	if cate.Id > 0 {
		post.CateId = cate.Id
	}

	posts, err := a.artRepo.PostGetList(post, page, num, artdate, req.Keyword)
	if err != nil {
		return response.Fail(c, 500, err)
	}

	h := gin.H{}
	h["list"] = posts

	builder := a.db.Model(post)

	if artdate != "" {
		builder.Where("DATE_FORMAT(created_at,'%Y-%m') = ?", artdate)
	}

	if req.Keyword != "" {
		builder.Where("title like ?", "%"+req.Keyword+"%")
	}

	total, err := builder.Count()
	pager := pagination.New(int(total), num, req.Page, 2)
	h["pageTotal"] = pager.TotalPages()

	if err != nil {
		return response.Fail(c, 500, err)
	}
	return response.Success(c, h)
}

func (a *Article) PrevNext(c *gee.Context) gee.Response {
	req := &struct {
		Id int `json:"id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(req); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	h := gin.H{}
	h["prev"] = map[string]interface{}{}
	h["next"] = map[string]interface{}{}

	prev, err := a.artRepo.PostPrev(req.Id)
	if err == nil {
		h["prev"] = gin.H{
			"id":    prev.Id,
			"title": prev.Title,
		}
	}
	next, err := a.artRepo.PostNext(req.Id)
	if err == nil {
		h["next"] = gin.H{
			"id":    next.Id,
			"title": next.Title,
		}
	}
	return response.Success(c, h)
}

func (a *Article) Detail(c *gee.Context) gee.Response {
	req := &struct {
		Id  int    `json:"id"`
		Url string `json:"url"`
	}{}
	if err := c.ShouldBindJSON(req); err != nil {
		return response.Fail(c, 201, "参数错误")
	}

	if req.Id < 1 && req.Url == "" {
		return response.Fail(c, 201, "参数错误")
	}

	post, err := a.artRepo.GetUserPost(req.Id, req.Url)
	if err != nil {
		return response.Fail(c, 202, "您访问的文章不存在或已经删除！")
	}

	return response.Success(c, post)
}

func (a *Article) Feed(c *gee.Context) gee.Response {
	now := time.Now()
	options, err := a.settingRepo.GetOptions()
	if err != nil {
		return response.Fail(c, 202, err)
	}

	feed := &feeds.Feed{
		Title:       options["site_name"],
		Link:        &feeds.Link{Href: "https://fifsky.com"},
		Description: options["site_desc"],
		Author:      &feeds.Author{Name: "fifsky", Email: "fifsky@gmail.com"},
		Created:     now,
	}

	cid := convert.StrTo(c.DefaultQuery("cid", "0")).MustInt()

	post := &model.Posts{}
	if cid > 0 {
		post.CateId = cid
	}

	posts, err := a.artRepo.PostGetList(post, 1, 10, "", "")

	if err != nil {
		return response.Fail(c, 500, err)
	}

	for _, v := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://fifsky.com/article/%d", v.Id)},
			Description: v.Content,
			Author:      &feeds.Author{Name: v.User.NickName, Email: "fifsky@gmail.com"},
			Created:     now,
		})
	}

	err = feed.WriteAtom(c.Writer)
	if err != nil {
		return response.Fail(c, 500, err)
	}
	return nil
}

func (a *Article) Post(c *gee.Context) gee.Response {
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
		if _, err := a.db.Model(post).Update(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "更新文章失败")
		}
	} else {
		if _, err := a.db.Model(post).Create(); err != nil {
			logger.Error(err)
			return response.Fail(c, 201, "发表文章失败")
		}
	}

	return response.Success(c, post)
}

func (a *Article) Delete(c *gee.Context) gee.Response {
	p := &struct {
		Id int `json:"id" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		return response.Fail(c, 201, "参数错误:"+err.Error())
	}

	post := &model.Posts{Id: p.Id, Status: 2}
	if _, err := a.db.Model(post).Update(); err != nil {
		logger.Error(err)
		return response.Fail(c, 201, "删除失败")
	}
	return response.Success(c, nil)
}

func (a *Article) Upload(c *gee.Context) gee.Response {
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
	_, err = file.Seek(0, 0)
	if err != nil {
		return c.JSON(gin.H{
			"errno": 203,
		})
	}

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
