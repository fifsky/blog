package admin

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"app/config"
	"app/pkg/errors"
	"app/pkg/ossutil"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/server/response"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.ArticleServiceHTTPServer = (*Article)(nil)

type Article struct {
	store *store.Store
	conf  *config.Config
	upl   ossutil.Uploader
}

func NewArticle(s *store.Store, conf *config.Config) *Article {
	return &Article{
		store: s,
		conf:  conf,
		upl:   ossutil.NewAliyunUploader(conf),
	}
}

func (a *Article) Create(ctx context.Context, req *adminv1.ArticleCreateRequest) (*types.IDResponse, error) {
	loginUser := GetLoginUser(ctx)
	now := time.Now()
	status := model.PostStatusActive
	if req.GetStatus() != "" {
		status = model.PostStatus(req.GetStatus())
	}
	c := &model.Post{
		CateId:    int(req.GetCateId()),
		Type:      int(req.GetType()),
		UserId:    loginUser.Id,
		Title:     req.GetTitle(),
		Url:       req.GetUrl(),
		Content:   req.GetContent(),
		Tags:      model.Tags(req.GetTags()),
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := a.store.CreatePost(ctx, c)
	if err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(lastId)}.Build(), nil
}

func (a *Article) Update(ctx context.Context, req *adminv1.ArticleUpdateRequest) (*types.IDResponse, error) {
	now := time.Now()
	u := &model.UpdatePost{Id: int(req.GetId())}
	if req.GetCateId() > 0 {
		u.CateId = new(int(req.GetCateId()))
	}
	if req.GetType() > 0 {
		u.Type = new(int(req.GetType()))
	}
	if req.GetStatus() != "" {
		u.Status = new(model.PostStatus(req.GetStatus()))
	}
	if req.GetTitle() != "" {
		u.Title = new(req.GetTitle())
	}
	// 始终更新 url 字段，允许清空自定义路径
	u.Url = new(req.GetUrl())
	if req.GetContent() != "" {
		u.Content = new(req.GetContent())
	}
	if req.GetTags() != nil {
		u.Tags = new(model.Tags(req.GetTags()))
	}
	u.UpdatedAt = &now
	if err := a.store.UpdatePost(ctx, u); err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: req.GetId()}.Build(), nil
}

func (a *Article) Delete(ctx context.Context, req *adminv1.ArticleDeleteRequest) (*emptypb.Empty, error) {
	ids := lo.FilterMap(req.GetIds(), func(id int32, _ int) (int, bool) {
		return int(id), id > 0
	})
	if len(ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := a.store.SoftDeletePost(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (a *Article) Restore(ctx context.Context, req *adminv1.ArticleRestoreRequest) (*types.IDResponse, error) {
	ids := lo.FilterMap(req.GetIds(), func(id int32, _ int) (int, bool) {
		return int(id), id > 0
	})
	if len(ids) == 0 {
		return types.IDResponse_builder{}.Build(), nil
	}

	// 批量恢复，逐个执行
	for _, id := range ids {
		if err := a.store.RestorePost(ctx, id); err != nil {
			return nil, err
		}
	}
	return types.IDResponse_builder{Id: int32(ids[0])}.Build(), nil
}

func (a *Article) Destroy(ctx context.Context, req *adminv1.ArticleDestroyRequest) (*emptypb.Empty, error) {
	ids := lo.FilterMap(req.GetIds(), func(id int32, _ int) (int, bool) {
		return int(id), id > 0
	})
	if len(ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := a.store.DestroyPost(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (a *Article) Detail(ctx context.Context, req *adminv1.ArticleDetailRequest) (*adminv1.ArticleItem, error) {
	post, err := a.store.GetPost(ctx, int(req.GetId()), "")
	if err != nil {
		return nil, err
	}

	item := adminv1.ArticleItem_builder{Id: int32(post.Id),
		CateId:    int32(post.CateId),
		Type:      int32(post.Type),
		UserId:    int32(post.UserId),
		Title:     post.Title,
		Url:       post.Url,
		Content:   post.Content,
		Tags:      []string(post.Tags),
		Status:    string(post.Status),
		ViewNum:   int32(post.ViewNum),
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

func (a *Article) List(ctx context.Context, req *adminv1.ArticleListRequest) (*adminv1.ArticleListResponse, error) {
	page := 1
	if req.GetPage() > 0 {
		page = int(req.GetPage())
	}
	num := 20
	posts, err := a.store.ListPostForAdmin(ctx, &model.Post{
		Type:   int(req.GetType()),
		Status: model.PostStatus(req.GetStatus()),
	}, page, num, req.GetKeyword())
	if err != nil {
		return nil, err
	}
	uids := lo.Map(posts, func(p model.Post, _ int) int { return p.UserId })
	cids := lo.Map(posts, func(p model.Post, _ int) int { return p.CateId })
	um, err := a.store.GetUserByIds(ctx, lo.Uniq(uids))
	if err != nil {
		return nil, err
	}
	cm, err := a.store.GetCatesByIds(ctx, lo.Uniq(cids))
	if err != nil {
		return nil, err
	}
	items := lo.Map(posts, func(p model.Post, _ int) *adminv1.ArticleItem {
		item := adminv1.ArticleItem_builder{Id: int32(p.Id),
			CateId:    int32(p.CateId),
			Type:      int32(p.Type),
			UserId:    int32(p.UserId),
			Title:     p.Title,
			Url:       p.Url,
			Content:   p.Content,
			Tags:      []string(p.Tags),
			Status:    string(p.Status),
			ViewNum:   int32(p.ViewNum),
			CreatedAt: p.CreatedAt.Format(time.DateTime),
			UpdatedAt: p.UpdatedAt.Format(time.DateTime)}.Build()

		if u, ok := um[p.UserId]; ok {
			item.SetUser(types.UserSummary_builder{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}.Build())
		}
		if c, ok := cm[p.CateId]; ok {
			item.SetCate(types.CateSummary_builder{Id: int32(c.Id), Name: c.Name, Domain: c.Domain}.Build())
		}
		return item
	})
	total, err := a.store.CountPostsForAdmin(ctx, &model.Post{
		Type:   int(req.GetType()),
		Status: model.PostStatus(req.GetStatus()),
	}, req.GetKeyword())
	if err != nil {
		return nil, err
	}
	return adminv1.ArticleListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

// Upload 上传接口（仅管理员）
func (a *Article) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Fail(w, errors.BadRequest("UPLOAD_FILE_ERROR", err.Error()))
		return
	}
	day := time.Now().Format("20060102")
	filename := "upload/" + day + "/" + md5File(file) + ".png"
	_, err = file.Seek(0, 0)
	if err != nil {
		response.Fail(w, errors.BadRequest("UPLOAD_FILE_ERROR", err.Error()))
		return
	}
	err = a.upl.Put(r.Context(), filename, file)
	if err != nil {
		response.Fail(w, errors.BadRequest("UPLOAD_FILE_ERROR", err.Error()))
		return
	}
	response.Success(w, map[string]any{
		"url": "https://static.fifsky.com/" + filename,
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
