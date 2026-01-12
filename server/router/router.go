package router

import (
	"log/slog"
	"net/http"
	"path/filepath"

	"app/config"
	"app/contract"
	"app/middleware"
	"app/pkg/logger"
	"app/pkg/logger/sloghttp"
	adminv1 "app/proto/gen/admin/v1"
	apiv1 "app/proto/gen/api/v1"
	"app/response"
	"app/service"
	adminsvc "app/service/admin"
	"app/store"
)

type Router struct {
	service *service.Service
	admin   *adminsvc.Service
	conf    *config.Config
	store   *store.Store
}

func New(service *service.Service, conf *config.Config, s *store.Store) Router {
	return Router{
		service: service,
		admin:   adminsvc.New(s, conf),
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
		response.Fail(w, 404, "接口不存在")
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
			sloghttp.IgnorePath("/api/admin/upload"),
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

	// 统一处理所有 /api/ 路径的预检请求，确保中间件设置CORS响应头
	api.HandleFunc("OPTIONS /api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.HandleFunc("GET /api/weixin/notify", r.service.Weixin.Notify)
	api.HandleFunc("POST /api/weixin/notify", r.service.Weixin.Notify)

	codec := contract.NewCodec()
	apiv1.RegisterArticleServiceHTTPServer(api, codec, r.service.Article)
	apiv1.RegisterMoodServiceHTTPServer(api, codec, r.service.Mood)
	apiv1.RegisterCateServiceHTTPServer(api, codec, r.service.Cate)
	apiv1.RegisterLinkServiceHTTPServer(api, codec, r.service.Link)
	apiv1.RegisterRemindServiceHTTPServer(api, codec, r.service.Remind)
	apiv1.RegisterUserServiceHTTPServer(api, codec, r.service.User)
	apiv1.RegisterSettingServiceHTTPServer(api, codec, r.service.Setting)
	apiv1.RegisterWeixinServiceHTTPServer(api, codec, r.service.Weixin)

	adminAuth := api.Use(middleware.NewAuthLogin(r.store, r.conf))
	adminAuth.HandleFunc("POST /api/admin/upload", r.admin.Article.Upload)
	adminv1.RegisterArticleServiceHTTPServer(adminAuth, codec, r.admin.Article)
	adminv1.RegisterMoodServiceHTTPServer(adminAuth, codec, r.admin.Mood)
	adminv1.RegisterCateServiceHTTPServer(adminAuth, codec, r.admin.Cate)
	adminv1.RegisterLinkServiceHTTPServer(adminAuth, codec, r.admin.Link)
	adminv1.RegisterRemindServiceHTTPServer(adminAuth, codec, r.admin.Remind)
	adminv1.RegisterUserServiceHTTPServer(adminAuth, codec, r.admin.User)
	adminv1.RegisterSettingServiceHTTPServer(adminAuth, codec, r.admin.Setting)

	return &NotFoundHandler{mux: mux.ServeMux}
}
