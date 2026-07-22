package router

import (
	"log/slog"
	"net/http"

	"app/config"
	"app/pkg/errors"
	adminv1 "app/proto/gen/admin/v1"
	apiv1 "app/proto/gen/api/v1"
	"app/server/codec"
	"app/server/middleware"
	"app/server/response"
	adminsvc "app/service/admin"
	"app/service/mcptool"
	"app/service/openapi"
	"app/store"

	"github.com/goapt/logger/sloghttp"
)

type Router struct {
	service      *openapi.Service
	admin        *adminsvc.Service
	conf         *config.Config
	store        *store.Store
	accessLogger *slog.Logger
}

func New(apiService *openapi.Service, adminService *adminsvc.Service, conf *config.Config, s *store.Store, accessLogger *slog.Logger) Router {
	return Router{
		service:      apiService,
		admin:        adminService,
		conf:         conf,
		store:        s,
		accessLogger: accessLogger,
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
		Level:              slog.LevelInfo,
		WithUserAgent:      true,
		WithRequestHeader:  true,
		WithRequestBody:    true,
		WithResponseHeader: true,
		// WithResponseBody: true,
		Filters: []sloghttp.Filter{
			sloghttp.IgnorePath("/blog/admin/upload", "/blog/china_map", "/blog/admin/ai/chat", "/blog/admin/ai/tags"),
		},
	}

	mux := NewServeMux()
	api := mux.Use(middleware.NewRecover, sloghttp.NewMiddleware(r.accessLogger, conf), middleware.NewHeader, middleware.NewCors)

	// 统一处理所有 /blog/ 路径的预检请求，确保中间件设置CORS响应头
	api.HandleFunc("OPTIONS /blog/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	codec := codec.NewCodec()
	apiv1.RegisterArticleServiceHTTPServer(api, codec, r.service.Article)
	apiv1.RegisterMoodServiceHTTPServer(api, codec, r.service.Mood)
	apiv1.RegisterCateServiceHTTPServer(api, codec, r.service.Cate)
	apiv1.RegisterLinkServiceHTTPServer(api, codec, r.service.Link)
	apiv1.RegisterUserServiceHTTPServer(api, codec, r.service.User)
	apiv1.RegisterSettingServiceHTTPServer(api, codec, r.service.Setting)
	apiv1.RegisterTravelServiceHTTPServer(api, codec, r.service.Travel)
	apiv1.RegisterMiniAppServiceHTTPServer(api, codec, r.service.MiniApp)
	apiv1.RegisterGeoServiceHTTPServer(api, codec, r.service.Geo)
	apiv1.RegisterGuestbookServiceHTTPServer(api, codec, r.service.Guestbook)
	apiv1.RegisterCommentServiceHTTPServer(api, codec, r.service.Comment)

	mcpAuth := api.Use(middleware.NewToken(r.conf.Common.MCPToken))
	mcpRemindHandler := mcptool.NewRemindHandler(r.store)
	mcpAuth.Handle("POST /blog/mcp/remind", mcpRemindHandler)
	mcpAuth.Handle("GET /blog/mcp/remind", mcpRemindHandler)

	mcpMoodHandler := mcptool.NewMoodHandler(r.store)
	mcpAuth.Handle("POST /blog/mcp/mood", mcpMoodHandler)
	mcpAuth.Handle("GET /blog/mcp/mood", mcpMoodHandler)

	adminAuth := api.Use(middleware.NewAuthLogin(r.store, r.conf.Common.TokenSecret))
	adminAuth.HandleFunc("POST /blog/admin/upload", r.admin.Article.Upload)
	adminv1.RegisterArticleServiceHTTPServer(adminAuth, codec, r.admin.Article)
	adminv1.RegisterMoodServiceHTTPServer(adminAuth, codec, r.admin.Mood)
	adminv1.RegisterCateServiceHTTPServer(adminAuth, codec, r.admin.Cate)
	adminv1.RegisterLinkServiceHTTPServer(adminAuth, codec, r.admin.Link)
	adminv1.RegisterRemindServiceHTTPServer(adminAuth, codec, r.admin.Remind)
	adminv1.RegisterUserServiceHTTPServer(adminAuth, codec, r.admin.User)
	adminv1.RegisterSettingServiceHTTPServer(adminAuth, codec, r.admin.Setting)
	adminv1.RegisterOSSServiceHTTPServer(adminAuth, codec, r.admin.OSS)
	adminv1.RegisterAIServiceHTTPServer(adminAuth, codec, r.admin.AI)
	adminv1.RegisterClawBotServiceHTTPServer(adminAuth, codec, r.admin.ClawBot)
	adminv1.RegisterGuestbookServiceHTTPServer(adminAuth, codec, r.admin.Guestbook)
	adminv1.RegisterCommentServiceHTTPServer(adminAuth, codec, r.admin.Comment)
	adminv1.RegisterFootprintServiceHTTPServer(adminAuth, codec, r.admin.Footprint)

	// AI chat endpoint (SSE streaming)
	adminAuth.HandleFunc("POST /blog/admin/ai/chat", r.admin.AI.Chat)

	return &NotFoundHandler{mux: mux.ServeMux}
}
