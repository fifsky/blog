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

func (a *Article) Create(ctx context.Context, req *adminv1.ArticleCreateRequest) (*adminv1.IDResponse, error) {
	loginUser := GetLoginUser(ctx)
	now := time.Now()
	c := &model.Post{
		CateId:    int(req.CateId),
		Type:      int(req.Type),
		UserId:    loginUser.Id,
		Title:     req.Title,
		Url:       req.Url,
		Content:   req.Content,
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	lastId, err := a.store.CreatePost(ctx, c)
	if err != nil {
		return nil, err
	}
	return &adminv1.IDResponse{Id: int32(lastId)}, nil
}

func (a *Article) Update(ctx context.Context, req *adminv1.ArticleUpdateRequest) (*adminv1.IDResponse, error) {
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
	return &adminv1.IDResponse{Id: req.Id}, nil
}

func (a *Article) Delete(ctx context.Context, req *adminv1.IDRequest) (*emptypb.Empty, error) {
	if err := a.store.SoftDeletePost(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
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
