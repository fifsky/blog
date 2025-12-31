package router

import (
	"net/http"

	"app/config"
	"app/handler"
	"app/middleware"
	"app/response"
	"app/service"
	"app/store"
)

type Router struct {
	handler *handler.Handler
	service *service.Service
	conf    *config.Config
	store   *store.Store
}

func New(handler *handler.Handler, service *service.Service, conf *config.Config, s *store.Store) Router {
	return Router{
		handler: handler,
		service: service,
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
	mux := NewServeMux()
	api := mux.Use(middleware.NewRecover, middleware.NewCors)

	// 统一处理所有 /api/ 路径的预检请求，确保中间件设置CORS响应头
	api.HandleFunc("OPTIONS /api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.HandleFunc("POST /api/login", r.handler.User.Login)
	api.HandleFunc("POST /api/mood/list", r.handler.Mood.List)
	api.HandleFunc("POST /api/cate/all", r.handler.Cate.All)
	api.HandleFunc("POST /api/article/archive", r.handler.Article.Archive)
	api.HandleFunc("POST /api/article/prevnext", r.handler.Article.PrevNext)
	api.HandleFunc("POST /api/article/list", r.handler.Article.List)
	api.HandleFunc("POST /api/article/detail", r.handler.Article.Detail)
	api.HandleFunc("POST /api/link/all", r.handler.Link.All)
	api.HandleFunc("POST /api/setting", r.handler.Setting.Get)
	api.HandleFunc("GET /feed.xml", r.handler.Article.Feed)

	remindAuth := api.Use(middleware.NewRemindAuth(r.store, r.conf))
	remindAuth.HandleFunc("GET /api/remind/change", r.handler.Remind.Change)
	remindAuth.HandleFunc("GET /api/remind/delay", r.handler.Remind.Delay)

	adminAuth := api.Use(middleware.NewAuthLogin(r.store, r.conf))
	adminAuth.HandleFunc("POST /api/admin/loginuser", r.handler.User.LoginUser)
	adminAuth.HandleFunc("POST /api/admin/article/create", r.handler.Article.Create)
	adminAuth.HandleFunc("POST /api/admin/article/update", r.handler.Article.Update)
	adminAuth.HandleFunc("POST /api/admin/article/delete", r.handler.Article.Delete)
	adminAuth.HandleFunc("POST /api/admin/upload", r.handler.Article.Upload)
	adminAuth.HandleFunc("POST /api/admin/setting/update", r.handler.Setting.Update)
	adminAuth.HandleFunc("POST /api/admin/mood/create", r.handler.Mood.Create)
	adminAuth.HandleFunc("POST /api/admin/mood/update", r.handler.Mood.Update)
	adminAuth.HandleFunc("POST /api/admin/mood/delete", r.handler.Mood.Delete)
	adminAuth.HandleFunc("POST /api/admin/cate/create", r.handler.Cate.Create)
	adminAuth.HandleFunc("POST /api/admin/cate/update", r.handler.Cate.Update)
	adminAuth.HandleFunc("POST /api/admin/cate/list", r.handler.Cate.List)
	adminAuth.HandleFunc("POST /api/admin/cate/delete", r.handler.Cate.Delete)
	adminAuth.HandleFunc("POST /api/admin/link/create", r.handler.Link.Create)
	adminAuth.HandleFunc("POST /api/admin/link/update", r.handler.Link.Update)
	adminAuth.HandleFunc("POST /api/admin/link/list", r.handler.Link.List)
	adminAuth.HandleFunc("POST /api/admin/link/delete", r.handler.Link.Delete)
	adminAuth.HandleFunc("POST /api/admin/remind/create", r.handler.Remind.Create)
	adminAuth.HandleFunc("POST /api/admin/remind/update", r.handler.Remind.Update)
	adminAuth.HandleFunc("POST /api/admin/remind/list", r.handler.Remind.List)
	adminAuth.HandleFunc("POST /api/admin/remind/delete", r.handler.Remind.Delete)
	adminAuth.HandleFunc("POST /api/admin/user/get", service.Wrap(r.service.User.Get))
	adminAuth.HandleFunc("POST /api/admin/user/create", r.handler.User.Create)
	adminAuth.HandleFunc("POST /api/admin/user/update", r.handler.User.Update)
	adminAuth.HandleFunc("POST /api/admin/user/list", r.handler.User.List)
	adminAuth.HandleFunc("POST /api/admin/user/status", r.handler.User.Status)

	return &NotFoundHandler{mux: mux.ServeMux}
}
