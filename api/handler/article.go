package handler

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"app/config"
	"app/model"
	"app/pkg/ossutil"
	"app/response"
	"app/store"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
)

type Article struct {
	store      *store.Store
	httpClient *http.Client
	conf       *config.Config
	uploader   ossutil.Uploader
}

func NewArticle(s *store.Store, conf *config.Config) *Article {
	return &Article{
		store:      s,
		httpClient: http.DefaultClient,
		conf:       conf,
		uploader:   ossutil.NewAliyunUploader(conf, http.DefaultClient),
	}
}

func (a *Article) Archive(w http.ResponseWriter, r *http.Request) {
	archives, err := a.store.PostArchive(r.Context())
	if err != nil {
		response.Fail(w, 500, err)
		return
	}
	data := make([]DateArchiveItem, 0)

	for _, v := range archives {
		data = append(data, DateArchiveItem{
			Url:     "/date/" + v.Ym,
			Content: fmt.Sprintf("%s(%s)", v.Ym, v.Total),
		})
	}

	response.Success(w, data)
}

func (a *Article) List(w http.ResponseWriter, r *http.Request) {
	req, err := decode[ArticleListRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	options, err := a.store.GetOptions(r.Context())
	if err != nil {
		response.Fail(w, 202, err)
		return
	}

	n, err := strconv.ParseInt(options["post_num"], 10, 0)
	if err != nil {
		n = 10
	}

	num := max(int(n), 1)

	cateId := 0
	if req.Domain != "" {
		cate, err := a.store.GetCateByDomain(r.Context(), req.Domain)
		if err != nil {
			response.Fail(w, 202, err)
			return
		}
		cateId = cate.Id
	}

	artdate := ""

	if req.Year != "" && req.Month != "" {
		artdate = req.Year + "-" + req.Month
	}

	page := max(req.Page, 1)

	post := &model.Post{
		CateId: cateId,
	}

	posts, err := a.store.ListPost(r.Context(), post, page, num, artdate, req.Keyword)
	if err != nil {
		response.Fail(w, 500, err)
		return
	}

	items := make([]ArticleItem, 0, len(posts))
	uids := lo.Map(posts, func(item model.Post, index int) int {
		return item.UserId
	})
	cids := lo.Map(posts, func(item model.Post, index int) int {
		return item.CateId
	})

	um, err := a.store.GetUserByIds(r.Context(), lo.Uniq(uids))
	if err != nil {
		response.Fail(w, 500, err)
		return
	}
	cm, err := a.store.GetCatesByIds(r.Context(), lo.Uniq(cids))
	if err != nil {
		response.Fail(w, 500, err)
		return
	}
	for _, p := range posts {
		item := ArticleItem{
			Id:        p.Id,
			CateId:    p.CateId,
			Type:      p.Type,
			UserId:    p.UserId,
			Title:     p.Title,
			Url:       p.Url,
			Content:   p.Content,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
		if u, ok := um[p.UserId]; ok {
			item.User = &UserSummary{Id: u.Id, Name: u.Name, NickName: u.NickName}
		}
		if c, ok := cm[p.CateId]; ok {
			item.Cate = &CateSummary{Id: c.Id, Name: c.Name, Domain: c.Domain}
		}
		items = append(items, item)
	}

	total, err := a.store.CountPosts(r.Context(), post, artdate, req.Keyword)
	resp := ArticleListResponse{
		List:      items,
		PageTotal: totalPages(total, num),
	}

	if err != nil {
		response.Fail(w, 500, err)
		return
	}
	response.Success(w, resp)
}

func (a *Article) PrevNext(w http.ResponseWriter, r *http.Request) {
	req, err := decode[PrevNextRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	resp := PrevNextResponse{
		Prev: PrevNextItem{},
		Next: PrevNextItem{},
	}

	prev, err := a.store.PrevPost(r.Context(), req.Id)
	if err == nil {
		resp.Prev = PrevNextItem{Id: prev.Id, Title: prev.Title}
	}
	next, err := a.store.NextPost(r.Context(), req.Id)
	if err == nil {
		resp.Next = PrevNextItem{Id: next.Id, Title: next.Title}
	}
	response.Success(w, resp)
}

func (a *Article) Detail(w http.ResponseWriter, r *http.Request) {
	req, err := decode[ArticleDetailRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误")
		return
	}

	post, err := a.store.GetPost(r.Context(), req.Id, req.Url)
	if err != nil {
		response.Fail(w, 202, "您访问的文章不存在或已经删除！")
		return
	}

	item := ArticleItem{
		Id:        post.Id,
		CateId:    post.CateId,
		Type:      post.Type,
		UserId:    post.UserId,
		Title:     post.Title,
		Url:       post.Url,
		Content:   post.Content,
		Status:    post.Status,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
	u, err := a.store.GetUser(r.Context(), post.UserId)
	if err != nil {
		response.Fail(w, 201, fmt.Sprintf("用户信息获取失败：%s", err))
		return
	}
	item.User = &UserSummary{Id: u.Id, Name: u.Name, NickName: u.NickName}
	c, err := a.store.GetCate(r.Context(), post.CateId)
	if err != nil {
		response.Fail(w, 201, fmt.Sprintf("文章分类获取失败：%s", err))
		return
	}
	item.Cate = &CateSummary{Id: c.Id, Name: c.Name, Domain: c.Domain}
	response.Success(w, item)
}

func (a *Article) Feed(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	options, err := a.store.GetOptions(r.Context())
	if err != nil {
		response.Fail(w, 202, err)
		return
	}

	feed := &feeds.Feed{
		Title:       options["site_name"],
		Link:        &feeds.Link{Href: "https://fifsky.com"},
		Description: options["site_desc"],
		Author:      &feeds.Author{Name: "fifsky", Email: "fifsky@gmail.com"},
		Created:     now,
	}

	cidStr := r.URL.Query().Get("cid")
	cid := 0
	if cidStr != "" {
		if v, e := strconv.Atoi(cidStr); e == nil {
			cid = v
		}
	}

	post := &model.Post{}
	if cid > 0 {
		post.CateId = cid
	}

	posts, err := a.store.ListPost(r.Context(), post, 1, 10, "", "")

	if err != nil {
		response.Fail(w, 500, err)
		return
	}

	uids := lo.Map(posts, func(item model.Post, index int) int {
		return item.UserId
	})

	um, _ := a.store.GetUserByIds(r.Context(), lo.Uniq(uids))

	for _, v := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://fifsky.com/article/%d", v.Id)},
			Description: v.Content,
			Author: &feeds.Author{Name: func() string {
				if u, ok := um[v.UserId]; ok {
					return u.NickName
				}
				return ""
			}(), Email: "fifsky@gmail.com"},
			Created: now,
		})
	}

	err = feed.WriteAtom(w)
	if err != nil {
		response.Fail(w, 500, err)
		return
	}
	// 无需返回JSON，直接输出Atom XML
}

func (a *Article) Post(w http.ResponseWriter, r *http.Request) {
	in, err := decode[ArticlePostRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	loginUser := getLoginUser(r.Context())
	if loginUser == nil {
		response.Fail(w, 201, "请先登录")
		return
	}

	now := time.Now()

	if in.Id > 0 {
		u := &model.UpdatePost{
			Id: in.Id,
		}
		if in.CateId > 0 {
			u.CateId = &in.CateId
		}
		if in.Type > 0 {
			u.Type = &in.Type
		}
		if in.Title != "" {
			u.Title = &in.Title
		}
		if in.Url != "" {
			u.Url = &in.Url
		}
		if in.Content != "" {
			u.Content = &in.Content
		}
		u.UpdatedAt = &now
		if err := a.store.UpdatePost(r.Context(), u); err != nil {
			response.Fail(w, 201, "更新文章失败")
			return
		}
	} else {
		c := &model.Post{
			CateId:    in.CateId,
			Type:      in.Type,
			UserId:    loginUser.Id,
			Title:     in.Title,
			Url:       in.Url,
			Content:   in.Content,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if _, err := a.store.CreatePost(r.Context(), c); err != nil {
			response.Fail(w, 201, "发表文章失败")
			return
		}
	}

	resp := ArticlePostResponse{
		Id:        in.Id,
		CateId:    in.CateId,
		Type:      in.Type,
		Title:     in.Title,
		Url:       in.Url,
		Content:   in.Content,
		Status:    1,
		CreatedAt: now.Format(time.DateTime),
		UpdatedAt: now.Format(time.DateTime),
	}
	response.Success(w, resp)
}

func (a *Article) Delete(w http.ResponseWriter, r *http.Request) {
	p, err := decode[IDRequest](r)
	if err != nil {
		response.Fail(w, 201, "参数错误:"+err.Error())
		return
	}

	if err := a.store.SoftDeletePost(r.Context(), p.Id); err != nil {
		response.Fail(w, 201, "删除失败")
		return
	}
	response.Success(w, nil)
}

func (a *Article) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad request\n" + err.Error()))
		return
	}
	if a.uploader == nil {
		a.uploader = ossutil.NewAliyunUploader(a.conf, a.httpClient)
	}

	day := time.Now().Format("20060102")

	filename := "upload/" + day + "/" + md5File(file) + ".png"
	_, err = file.Seek(0, 0)
	if err != nil {
		response.Upload(w, map[string]any{"errno": 203})
		return
	}

	err = a.uploader.Put(r.Context(), filename, file)
	if err != nil {
		response.Upload(w, map[string]any{"errno": 202})
		return
	}

	response.Upload(w, map[string]any{
		"errno": 0,
		"data": []string{
			"https://static.fifsky.com/" + filename + "!blog",
		},
	})
}
func md5File(file io.Reader) string {
	r := bufio.NewReader(file)
	md5h := md5.New()
	_, err := io.Copy(md5h, r)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(md5h.Sum(nil))
}
