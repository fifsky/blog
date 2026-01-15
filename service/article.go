package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"app/config"
	"app/model"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/proto/gen/types"
	"app/store"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
)

var _ apiv1.ArticleServiceServer = (*Article)(nil)

type Article struct {
	apiv1.UnimplementedArticleServiceServer
	store *store.Store
}

func NewArticle(s *store.Store, _ *config.Config) *Article {
	return &Article{
		store: s,
	}
}

func (a *Article) Archive(ctx context.Context, _ *emptypb.Empty) (*apiv1.ArchiveResponse, error) {
	archives, err := a.store.PostArchive(ctx)
	if err != nil {
		return nil, err
	}
	resp := &apiv1.ArchiveResponse{}
	for _, v := range archives {
		resp.List = append(resp.List, &apiv1.DateArchiveItem{
			Url:     "/date/" + v.Ym,
			Content: fmt.Sprintf("%s(%s)", v.Ym, v.Total),
		})
	}
	return resp, nil
}

func (a *Article) Calendar(ctx context.Context, req *apiv1.ArticleCalendarRequest) (*apiv1.ArticleCalendarResponse, error) {
	days, err := a.store.GetPostDaysInMonth(ctx, int(req.Year), int(req.Month))
	if err != nil {
		return nil, err
	}
	return &apiv1.ArticleCalendarResponse{
		Days: days,
	}, nil
}

func (a *Article) List(ctx context.Context, req *apiv1.ArticleListRequest) (*apiv1.ArticleListResponse, error) {
	options, err := a.store.GetOptions(ctx)
	if err != nil {
		return nil, err
	}
	n, err := strconv.ParseInt(options["post_num"], 10, 0)
	if err != nil {
		n = 10
	}
	num := max(int(n), 1)

	cateId := 0
	if req.Domain != "" {
		cate, err := a.store.GetCateByDomain(ctx, req.Domain)
		if err != nil {
			return nil, err
		}
		cateId = cate.Id
	}
	artdate := ""
	if req.Year != "" && req.Month != "" {
		artdate = req.Year + "-" + req.Month
		if req.Day != "" {
			artdate += "-" + req.Day
		}
	}
	page := max(int(req.Page), 1)
	posts, err := a.store.ListPost(ctx, &model.Post{
		CateId: cateId,
	}, page, num, artdate, req.Keyword)
	if err != nil {
		return nil, err
	}
	items := make([]*apiv1.ArticleItem, 0, len(posts))
	uids := lo.Map(posts, func(item model.Post, index int) int {
		return item.UserId
	})
	cids := lo.Map(posts, func(item model.Post, index int) int {
		return item.CateId
	})
	um, err := a.store.GetUserByIds(ctx, lo.Uniq(uids))
	if err != nil {
		return nil, err
	}
	cm, err := a.store.GetCatesByIds(ctx, lo.Uniq(cids))
	if err != nil {
		return nil, err
	}
	for _, p := range posts {
		item := &apiv1.ArticleItem{
			Id:        int32(p.Id),
			CateId:    int32(p.CateId),
			Type:      int32(p.Type),
			UserId:    int32(p.UserId),
			Title:     p.Title,
			Url:       p.Url,
			Content:   p.Content,
			Status:    int32(p.Status),
			CreatedAt: p.CreatedAt.Format(time.DateTime),
			UpdatedAt: p.UpdatedAt.Format(time.DateTime),
		}
		if u, ok := um[p.UserId]; ok {
			item.User = &types.UserSummary{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}
		}
		if c, ok := cm[p.CateId]; ok {
			item.Cate = &types.CateSummary{Id: int32(c.Id), Name: c.Name, Domain: c.Domain}
		}
		items = append(items, item)
	}
	total, err := a.store.CountPosts(ctx, &model.Post{CateId: cateId}, artdate, req.Keyword)
	if err != nil {
		return nil, err
	}
	return &apiv1.ArticleListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (a *Article) PrevNext(ctx context.Context, req *apiv1.PrevNextRequest) (*apiv1.PrevNextResponse, error) {
	resp := &apiv1.PrevNextResponse{}
	if prev, err := a.store.PrevPost(ctx, int(req.Id)); err == nil {
		resp.Prev = &apiv1.PrevNextItem{Id: int32(prev.Id), Title: prev.Title}
	}
	if next, err := a.store.NextPost(ctx, int(req.Id)); err == nil {
		resp.Next = &apiv1.PrevNextItem{Id: int32(next.Id), Title: next.Title}
	}
	return resp, nil
}

func (a *Article) Detail(ctx context.Context, req *apiv1.ArticleDetailRequest) (*apiv1.ArticleItem, error) {
	post, err := a.store.GetPost(ctx, int(req.Id), req.Url)
	if err != nil {
		return nil, errors.ErrSystem.WithCause(err)
	}

	if post.Status != 1 {
		return nil, errors.ErrArticleNotFound
	}

	item := &apiv1.ArticleItem{
		Id:        int32(post.Id),
		CateId:    int32(post.CateId),
		Type:      int32(post.Type),
		UserId:    int32(post.UserId),
		Title:     post.Title,
		Url:       post.Url,
		Content:   post.Content,
		Status:    int32(post.Status),
		CreatedAt: post.CreatedAt.Format(time.DateTime),
		UpdatedAt: post.UpdatedAt.Format(time.DateTime),
	}
	u, err := a.store.GetUser(ctx, post.UserId)
	if err == nil {
		item.User = &types.UserSummary{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}
	}
	c, err := a.store.GetCate(ctx, post.CateId)
	if err == nil {
		item.Cate = &types.CateSummary{Id: int32(c.Id), Name: c.Name, Domain: c.Domain}
	}
	return item, nil
}

// Feed 返回Atom XML
func (a *Article) Feed(ctx context.Context, _ *emptypb.Empty) (*httpbody.HttpBody, error) {
	now := time.Now()
	options, err := a.store.GetOptions(ctx)
	if err != nil {
		return nil, err
	}
	feed := &feeds.Feed{
		Title:       options["site_name"],
		Link:        &feeds.Link{Href: "https://fifsky.com"},
		Description: options["site_desc"],
		Author:      &feeds.Author{Name: "fifsky", Email: "fifsky@gmail.com"},
		Created:     now,
	}
	posts, err := a.store.ListPost(ctx, &model.Post{}, 1, 10, "", "")
	if err != nil {
		return nil, err
	}
	uids := lo.Map(posts, func(item model.Post, index int) int {
		return item.UserId
	})
	um, _ := a.store.GetUserByIds(ctx, lo.Uniq(uids))
	for _, v := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://fifsky.com/article/%d", v.Id)},
			Description: v.Content,
			Author: &feeds.Author{
				Name: func() string {
					if u, ok := um[v.UserId]; ok {
						return u.NickName
					}
					return ""
				}(),
				Email: "fifsky@gmail.com",
			},
			Created: now,
		})
	}
	atom, err := feed.ToAtom()
	if err != nil {
		return nil, err
	}
	return &httpbody.HttpBody{
		ContentType: "application/xml; charset=utf-8",
		Data:        []byte(atom),
	}, nil
}
