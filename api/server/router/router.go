package router

import (
	"net/http"

	"app/config"
	"app/handler"
	"app/middleware"
	"app/response"
	"app/store"
)

type Router struct {
	handler *handler.Handler
	conf    *config.Config
	store   *store.Store
}

func New(handler *handler.Handler, conf *config.Config, s *store.Store) Router {
	return Router{
		handler: handler,
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
	mux := http.NewServeMux()
	mid := Use(middleware.NewRecover, middleware.NewCors(r.conf))

	mux.Handle("POST /api/login", mid.HandlerFunc(r.handler.User.Login))
	mux.Handle("POST /api/mood/list", mid.HandlerFunc(r.handler.Mood.List))
	mux.Handle("POST /api/cate/all", mid.HandlerFunc(r.handler.Cate.All))
	mux.Handle("POST /api/article/archive", mid.HandlerFunc(r.handler.Article.Archive))
	mux.Handle("POST /api/article/prevnext", mid.HandlerFunc(r.handler.Article.PrevNext))
	mux.Handle("POST /api/article/list", mid.HandlerFunc(r.handler.Article.List))
	mux.Handle("POST /api/article/detail", mid.HandlerFunc(r.handler.Article.Detail))
	mux.Handle("POST /api/link/all", mid.HandlerFunc(r.handler.Link.All))
	// mux.Handle("POST /api/comment/new", mid.HandlerFunc(r.handler.Comment.Add)
	// mux.Handle("POST /api/comment/list", mid.HandlerFunc(r.handler.Comment.List)
	// mux.Handle("POST /api/comment/post", mid.HandlerFunc(r.handler.Comment.Post)

	remindAuth := mid.Append(middleware.NewRemindAuth(r.store, r.conf))
	mux.Handle("GET /api/remind/change", remindAuth.HandlerFunc(r.handler.Remind.Change))
	mux.Handle("GET /api/remind/delay", remindAuth.HandlerFunc(r.handler.Remind.Delay))
	mux.Handle("POST /api/setting", mid.HandlerFunc(r.handler.Setting.Get))
	mux.Handle("GET /feed.xml", mid.HandlerFunc(r.handler.Article.Feed))

	adminAuth := mid.Append(middleware.NewAuthLogin(r.store, r.conf))
	mux.Handle("POST /api/admin/loginuser", adminAuth.HandlerFunc(r.handler.User.LoginUser))
	mux.Handle("POST /api/admin/article/create", adminAuth.HandlerFunc(r.handler.Article.Create))
	mux.Handle("POST /api/admin/article/update", adminAuth.HandlerFunc(r.handler.Article.Update))
	mux.Handle("POST /api/admin/article/delete", adminAuth.HandlerFunc(r.handler.Article.Delete))
	mux.Handle("POST /api/admin/upload", adminAuth.HandlerFunc(r.handler.Article.Upload))
	mux.Handle("POST /api/admin/setting/update", adminAuth.HandlerFunc(r.handler.Setting.Update))
	// mux.Handle("POST /api/admin/comment/list", adminAuth.HandlerFunc(r.handler.Comment.AdminList))
	// mux.Handle("POST /api/admin/comment/delete", adminAuth.HandlerFunc(r.handler.Comment.Delete))
	mux.Handle("POST /api/admin/mood/create", adminAuth.HandlerFunc(r.handler.Mood.Create))
	mux.Handle("POST /api/admin/mood/update", adminAuth.HandlerFunc(r.handler.Mood.Update))
	mux.Handle("POST /api/admin/mood/delete", adminAuth.HandlerFunc(r.handler.Mood.Delete))
	mux.Handle("POST /api/admin/cate/create", adminAuth.HandlerFunc(r.handler.Cate.Create))
	mux.Handle("POST /api/admin/cate/update", adminAuth.HandlerFunc(r.handler.Cate.Update))
	mux.Handle("POST /api/admin/cate/list", adminAuth.HandlerFunc(r.handler.Cate.List))
	mux.Handle("POST /api/admin/cate/delete", adminAuth.HandlerFunc(r.handler.Cate.Delete))
	mux.Handle("POST /api/admin/link/create", adminAuth.HandlerFunc(r.handler.Link.Create))
	mux.Handle("POST /api/admin/link/update", adminAuth.HandlerFunc(r.handler.Link.Update))
	mux.Handle("POST /api/admin/link/list", adminAuth.HandlerFunc(r.handler.Link.List))
	mux.Handle("POST /api/admin/link/delete", adminAuth.HandlerFunc(r.handler.Link.Delete))
	mux.Handle("POST /api/admin/remind/create", adminAuth.HandlerFunc(r.handler.Remind.Create))
	mux.Handle("POST /api/admin/remind/update", adminAuth.HandlerFunc(r.handler.Remind.Update))
	mux.Handle("POST /api/admin/remind/list", adminAuth.HandlerFunc(r.handler.Remind.List))
	mux.Handle("POST /api/admin/remind/delete", adminAuth.HandlerFunc(r.handler.Remind.Delete))
	mux.Handle("POST /api/admin/user/get", adminAuth.HandlerFunc(r.handler.User.Get))
	mux.Handle("POST /api/admin/user/create", adminAuth.HandlerFunc(r.handler.User.Create))
	mux.Handle("POST /api/admin/user/update", adminAuth.HandlerFunc(r.handler.User.Update))
	mux.Handle("POST /api/admin/user/list", adminAuth.HandlerFunc(r.handler.User.List))
	mux.Handle("POST /api/admin/user/status", adminAuth.HandlerFunc(r.handler.User.Status))

	return &NotFoundHandler{mux: mux}
}
