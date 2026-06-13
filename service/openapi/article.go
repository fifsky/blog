package openapi

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"app/config"
	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
)

var _ apiv1.ArticleServiceHTTPServer = (*Article)(nil)

type Article struct {
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
	resp := apiv1.ArchiveResponse_builder{}.Build()
	for _, v := range archives {
		resp.SetList(append(resp.GetList(), apiv1.DateArchiveItem_builder{Url: "/date/" + v.Ym,
			Content: fmt.Sprintf("%s(%s)", v.Ym, v.Total)}.Build(),
		))
	}
	return resp, nil
}

func (a *Article) Calendar(ctx context.Context, req *apiv1.ArticleCalendarRequest) (*apiv1.ArticleCalendarResponse, error) {
	days, err := a.store.GetPostDaysInMonth(ctx, int(req.GetYear()), int(req.GetMonth()))
	if err != nil {
		return nil, err
	}
	return apiv1.ArticleCalendarResponse_builder{Days: days}.Build(),
		nil
}

func (a *Article) List(ctx context.Context, req *apiv1.ArticleListRequest) (*apiv1.ArticleListResponse, error) {
	options, err := a.store.GetOptions(ctx)
	if err != nil {
		return nil, err
	}
	var num int
	if req.GetPageSize() > 0 {
		num = int(req.GetPageSize())
	} else {
		n, err := strconv.ParseInt(options["post_num"], 10, 0)
		if err != nil {
			n = 10
		}
		num = max(int(n), 1)
	}

	cateId := 0
	if req.GetDomain() != "" {
		cate, err := a.store.GetCateByDomain(ctx, req.GetDomain())
		if err != nil {
			return nil, err
		}
		cateId = cate.Id
	}
	artdate := ""
	if req.GetYear() != "" && req.GetMonth() != "" {
		artdate = req.GetYear() + "-" + req.GetMonth()
		if req.GetDay() != "" {
			artdate += "-" + req.GetDay()
		}
	}
	page := max(int(req.GetPage()), 1)
	posts, err := a.store.ListPost(ctx, &model.Post{
		CateId: cateId,
		Type:   int(req.GetType()),
	}, page, num, artdate, req.GetKeyword(), req.GetTag())
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
		item := apiv1.ArticleItem_builder{Id: int32(p.Id),
			CateId:    int32(p.CateId),
			Type:      int32(p.Type),
			UserId:    int32(p.UserId),
			Title:     p.Title,
			Url:       p.Url,
			Content:   p.Content,
			Tags:      []string(p.Tags),
			Status:    int32(p.Status),
			CreatedAt: p.CreatedAt.Format(time.DateTime),
			UpdatedAt: p.UpdatedAt.Format(time.DateTime)}.Build()

		if u, ok := um[p.UserId]; ok {
			item.SetUser(types.UserSummary_builder{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}.Build())
		}
		if c, ok := cm[p.CateId]; ok {
			item.SetCate(types.CateSummary_builder{Id: int32(c.Id), Name: c.Name, Domain: c.Domain}.Build())
		}
		items = append(items, item)
	}
	total, err := a.store.CountPosts(ctx, &model.Post{CateId: cateId, Type: int(req.GetType())}, artdate, req.GetKeyword(), req.GetTag())
	if err != nil {
		return nil, err
	}
	return apiv1.ArticleListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (a *Article) PrevNext(ctx context.Context, req *apiv1.PrevNextRequest) (*apiv1.PrevNextResponse, error) {
	resp := apiv1.PrevNextResponse_builder{}.Build()
	if prev, err := a.store.PrevPost(ctx, int(req.GetId())); err == nil {
		resp.SetPrev(apiv1.PrevNextItem_builder{Id: int32(prev.Id), Title: prev.Title, Url: prev.Url}.Build())
	}
	if next, err := a.store.NextPost(ctx, int(req.GetId())); err == nil {
		resp.SetNext(apiv1.PrevNextItem_builder{Id: int32(next.Id), Title: next.Title, Url: next.Url}.Build())
	}
	return resp, nil
}

func (a *Article) Detail(ctx context.Context, req *apiv1.ArticleDetailRequest) (*apiv1.ArticleItem, error) {
	post, err := a.store.GetPost(ctx, int(req.GetId()), req.GetUrl())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrArticleNotFound
		}
		return nil, errors.ErrSystem.WithCause(err)
	}

	if post.Status != 1 {
		return nil, errors.ErrArticleNotFound
	}

	// Increment view count
	_ = a.store.IncrementPostViewNum(ctx, post.Id)

	item := apiv1.ArticleItem_builder{Id: int32(post.Id),
		CateId:    int32(post.CateId),
		Type:      int32(post.Type),
		UserId:    int32(post.UserId),
		Title:     post.Title,
		Url:       post.Url,
		Content:   post.Content,
		Tags:      []string(post.Tags),
		Status:    int32(post.Status),
		ViewNum:   int32(post.ViewNum + 1), // Return incremented value
		CreatedAt: post.CreatedAt.Format(time.DateTime),
		UpdatedAt: post.UpdatedAt.Format(time.DateTime)}.Build()

	u, err := a.store.GetUser(ctx, post.UserId)
	if err == nil {
		item.SetUser(types.UserSummary_builder{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}.Build())
	}
	c, err := a.store.GetCate(ctx, post.CateId)
	if err == nil {
		item.SetCate(types.CateSummary_builder{Id: int32(c.Id), Name: c.Name, Domain: c.Domain}.Build())
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
	posts, err := a.store.ListPost(ctx, &model.Post{}, 1, 10, "", "", "")
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
