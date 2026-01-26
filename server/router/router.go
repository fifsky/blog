package router

import (
	"log/slog"
	"net/http"
	"path/filepath"

	"app/config"
	"app/contract"
	"app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
	apiv1 "app/proto/gen/api/v1"
	"app/server/middleware"
	"app/server/response"
	adminsvc "app/service/admin"
	"app/service/openapi"
	"app/store"

	"github.com/goapt/logger"
	"github.com/goapt/logger/sloghttp"
)

type Router struct {
	service *openapi.Service
	admin   *adminsvc.Service
	conf    *config.Config
	store   *store.Store
}

func New(apiService *openapi.Service, adminService *adminsvc.Service, conf *config.Config, s *store.Store) Router {
	return Router{
		service: apiService,
		admin:   adminService,
		conf:    conf,
		store:   s,
	}
}

type NotFoundHandler struct {
	mux *http.ServeMux
}

func (n *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := n.mux.Handler(r)
	if pattern == "" || (pattern == "/" && r.URL.Path != "/") {
		response.Fail(w, errors.ErrApiNotFound)
		return
	}
	n.mux.ServeHTTP(w, r)
}

func (r *Router) Handler() http.Handler {
	conf := sloghttp.Config{
		DefaultLevel:       slog.LevelInfo,
		ClientErrorLevel:   slog.LevelWarn,
		ServerErrorLevel:   slog.LevelError,
		WithRequestID:      true,
		WithUserAgent:      true,
		WithRequestHeader:  true,
		WithRequestBody:    true,
		WithResponseHeader: true,
		// WithResponseBody: true,
		Filters: []sloghttp.Filter{
			sloghttp.IgnorePath("/blog/admin/upload", "/blog/china_map", "/blog/admin/ai/chat", "/blog/admin/ai/tags"),
		},
	}

	log := logger.New(&logger.Config{
		Mode:     logger.ModeFile,
		FileName: filepath.Join(r.conf.Common.StoragePath, "logs", "access.log"),
		MaxFiles: 3,
		Detail:   true,
	})

	mux := NewServeMux()
	api := mux.Use(middleware.NewRecover, sloghttp.NewMiddleware(log, conf), middleware.NewCors)

	// 统一处理所有 /blog/ 路径的预检请求，确保中间件设置CORS响应头
	api.HandleFunc("OPTIONS /blog/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.HandleFunc("GET /blog/weixin/notify", r.service.Weixin.Notify)
	api.HandleFunc("POST /blog/weixin/notify", r.service.Weixin.Notify)
	api.HandleFunc("GET /blog/region/nearest", r.service.Geo.Nearest)

	codec := contract.NewCodec()
	apiv1.RegisterArticleServiceHTTPServer(api, codec, r.service.Article)
	apiv1.RegisterMoodServiceHTTPServer(api, codec, r.service.Mood)
	apiv1.RegisterCateServiceHTTPServer(api, codec, r.service.Cate)
	apiv1.RegisterLinkServiceHTTPServer(api, codec, r.service.Link)
	apiv1.RegisterRemindServiceHTTPServer(api, codec, r.service.Remind)
	apiv1.RegisterUserServiceHTTPServer(api, codec, r.service.User)
	apiv1.RegisterSettingServiceHTTPServer(api, codec, r.service.Setting)
	apiv1.RegisterWeixinServiceHTTPServer(api, codec, r.service.Weixin)
	apiv1.RegisterTravelServiceHTTPServer(api, codec, r.service.Travel)
	apiv1.RegisterMiniAppServiceHTTPServer(api, codec, r.service.MiniApp)

	adminAuth := api.Use(middleware.NewAuthLogin(r.store, r.conf))
	adminAuth.HandleFunc("POST /blog/admin/upload", r.admin.Article.Upload)
	adminv1.RegisterArticleServiceHTTPServer(adminAuth, codec, r.admin.Article)
	adminv1.RegisterMoodServiceHTTPServer(adminAuth, codec, r.admin.Mood)
	adminv1.RegisterCateServiceHTTPServer(adminAuth, codec, r.admin.Cate)
	adminv1.RegisterLinkServiceHTTPServer(adminAuth, codec, r.admin.Link)
	adminv1.RegisterRemindServiceHTTPServer(adminAuth, codec, r.admin.Remind)
	adminv1.RegisterUserServiceHTTPServer(adminAuth, codec, r.admin.User)
	adminv1.RegisterSettingServiceHTTPServer(adminAuth, codec, r.admin.Setting)
	adminv1.RegisterPhotoServiceHTTPServer(adminAuth, codec, r.admin.Photo)
	adminv1.RegisterOSSServiceHTTPServer(adminAuth, codec, r.admin.OSS)
	adminv1.RegisterRegionServiceHTTPServer(adminAuth, codec, r.admin.Region)

	// AI chat endpoint (SSE streaming)
	adminAuth.HandleFunc("POST /blog/admin/ai/chat", r.admin.AI.Chat)
	adminAuth.HandleFunc("POST /blog/admin/ai/tags", r.admin.AI.GenerateTags)
	adminAuth.HandleFunc("POST /blog/admin/ai/remind/create", r.admin.AI.RemindSmartCreate)

	return &NotFoundHandler{mux: mux.ServeMux}
}
