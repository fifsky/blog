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
	"app/model"
	"app/pkg/ossutil"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/response"
	"app/store"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.ArticleServiceServer = (*Article)(nil)

type Article struct {
	adminv1.UnimplementedArticleServiceServer
	store *store.Store
	conf  *config.Config
	http  *http.Client
	upl   ossutil.Uploader
}

func NewArticle(s *store.Store, conf *config.Config) *Article {
	return &Article{
		store: s,
		conf:  conf,
		http:  http.DefaultClient,
		upl:   ossutil.NewAliyunUploader(conf, http.DefaultClient),
	}
}

func (a *Article) Create(ctx context.Context, req *adminv1.ArticleCreateRequest) (*types.IDResponse, error) {
	loginUser := GetLoginUser(ctx)
	now := time.Now()
	status := int32(1)
	if req.Status > 0 {
		status = req.Status
	}
	c := &model.Post{
		CateId:    int(req.CateId),
		Type:      int(req.Type),
		UserId:    loginUser.Id,
		Title:     req.Title,
		Url:       req.Url,
		Content:   req.Content,
		Status:    int(status),
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := a.store.CreatePost(ctx, c)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastId)}, nil
}

func (a *Article) Update(ctx context.Context, req *adminv1.ArticleUpdateRequest) (*types.IDResponse, error) {
	now := time.Now()
	u := &model.UpdatePost{Id: int(req.Id)}
	if req.CateId > 0 {
		v := int(req.CateId)
		u.CateId = &v
	}
	if req.Type > 0 {
		v := int(req.Type)
		u.Type = &v
	}
	if req.Title != "" {
		v := req.Title
		u.Title = &v
	}
	if req.Url != "" {
		v := req.Url
		u.Url = &v
	}
	if req.Content != "" {
		v := req.Content
		u.Content = &v
	}
	u.UpdatedAt = &now
	if err := a.store.UpdatePost(ctx, u); err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: req.Id}, nil
}

func (a *Article) Delete(ctx context.Context, req *adminv1.ArticleDeleteRequest) (*emptypb.Empty, error) {
	ids := make([]int, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id > 0 {
			ids = append(ids, int(id))
		}
	}
	if len(ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := a.store.SoftDeletePost(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (a *Article) Restore(ctx context.Context, req *adminv1.ArticleRestoreRequest) (*types.IDResponse, error) {
	ids := make([]int, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id > 0 {
			ids = append(ids, int(id))
		}
	}
	if len(ids) == 0 {
		return &types.IDResponse{}, nil
	}
	// 批量恢复，逐个执行
	for _, id := range ids {
		if err := a.store.RestorePost(ctx, id); err != nil {
			return nil, err
		}
	}
	return &types.IDResponse{Id: int32(ids[0])}, nil
}

func (a *Article) List(ctx context.Context, req *adminv1.ArticleListRequest) (*adminv1.ArticleListResponse, error) {
	page := 1
	if req.Page > 0 {
		page = int(req.Page)
	}
	num := 20
	posts, err := a.store.ListPostForAdmin(ctx, &model.Post{
		Type:   int(req.Type),
		Status: int(req.Status),
	}, page, num)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.ArticleItem, 0, len(posts))
	uids := make([]int, 0, len(posts))
	cids := make([]int, 0, len(posts))
	for _, p := range posts {
		uids = append(uids, p.UserId)
		cids = append(cids, p.CateId)
	}
	um, err := a.store.GetUserByIds(ctx, uids)
	if err != nil {
		return nil, err
	}
	cm, err := a.store.GetCatesByIds(ctx, cids)
	if err != nil {
		return nil, err
	}
	for _, p := range posts {
		item := &adminv1.ArticleItem{
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
	total, err := a.store.CountPostsForAdmin(ctx, &model.Post{
		Type:   int(req.Type),
		Status: int(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return &adminv1.ArticleListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

// Upload 上传接口（仅管理员）
func (a *Article) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Upload(w, map[string]any{"errno": 201, "message": err.Error()})
		return
	}
	day := time.Now().Format("20060102")
	filename := "upload/" + day + "/" + md5File(file) + ".png"
	_, err = file.Seek(0, 0)
	if err != nil {
		response.Upload(w, map[string]any{"errno": 203, "message": err.Error()})
		return
	}
	err = a.upl.Put(r.Context(), filename, file)
	if err != nil {
		response.Upload(w, map[string]any{"errno": 202, "message": err.Error()})
		return
	}
	response.Upload(w, map[string]any{
		"errno": 0,
		"data": map[string]any{
			"url": "https://static.fifsky.com/" + filename,
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
