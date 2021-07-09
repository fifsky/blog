package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"app/response"
	"github.com/goapt/gee"
	"github.com/google/wire"

	"app/handler"
	"app/middleware"
)

type Router struct {
	handler    *handler.Handler
	middleware *middleware.Middleware
}

func NewRouter(handler *handler.Handler, middleware *middleware.Middleware) Router {
	return Router{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *Router) Run(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r.route(r.handler, r.middleware),
	}
	log.Println("[HTTP] Server listen:" + addr)
	gee.RegisterShutDown(func(sig os.Signal) {
		ctxw, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Close()
		if err := srv.Shutdown(ctxw); err != nil {
			log.Fatal("HTTP Server Shutdown:", err)
		}
		log.Println("HTTP Server exiting")
	})

	// service connections
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP listen: %s\n", err)
	}
}

func (r *Router) route(handler *handler.Handler, middleware *middleware.Middleware) http.Handler {
	router := gee.New()
	// cors middleware
	router.Use(gee.HandlerFunc(r.middleware.Cors))
	// panic recover middleware
	router.Use(gee.HandlerFunc(r.middleware.Recover))
	// log middleware use for all handle
	router.Use(gee.HandlerFunc(r.middleware.AccessLog))

	router.NoRoute(func(c *gee.Context) gee.Response {
		return response.Fail(c, 404, "接口不存在")
	})

	router.POST("/api/dingmsg", handler.DingTalk.DingMsg)
	router.POST("/api/login", handler.User.Login)
	router.POST("/api/mood/list", handler.Mood.List)
	router.POST("/api/cate/all", handler.Cate.All)
	router.POST("/api/article/archive", handler.Article.Archive)
	router.POST("/api/article/prevnext", handler.Article.PrevNext)
	router.POST("/api/article/list", handler.Article.List)
	router.POST("/api/article/detail", handler.Article.Detail)
	router.POST("/api/link/all", handler.Link.All)
	// router.POST("/api/comment/new", handler.Comment.Add)
	// router.POST("/api/comment/list", handler.Comment.List)
	// router.POST("/api/comment/post", handler.Comment.Post)
	router.GET("/api/avatar", handler.Common.Avatar)
	router.GET("/api/remind/change", gee.HandlerFunc(middleware.RemindAuth), handler.Remind.Change)
	router.GET("/api/remind/delay", gee.HandlerFunc(middleware.RemindAuth), handler.Remind.Delay)
	router.POST("/api/setting", handler.Setting.Get)
	router.GET("/feed.xml", handler.Article.Feed)

	admin := router.Group("/api/admin")
	admin.Use(gee.HandlerFunc(middleware.AuthLogin))
	{
		admin.POST("/loginuser", handler.User.LoginUser)
		admin.POST("/article/post", handler.Article.Post)
		admin.POST("/article/delete", handler.Article.Delete)
		admin.POST("/upload", handler.Article.Upload)
		admin.POST("/setting/post", handler.Setting.Post)
		admin.POST("/comment/list", handler.Comment.AdminList)
		admin.POST("/comment/delete", handler.Comment.Delete)
		admin.POST("/mood/post", handler.Mood.Post)
		admin.POST("/mood/delete", handler.Mood.Delete)
		admin.POST("/cate/post", handler.Cate.Post)
		admin.POST("/cate/list", handler.Cate.List)
		admin.POST("/cate/delete", handler.Cate.Delete)
		admin.POST("/link/post", handler.Link.Post)
		admin.POST("/link/list", handler.Link.List)
		admin.POST("/link/delete", handler.Link.Delete)
		admin.POST("/remind/post", handler.Remind.Post)
		admin.POST("/remind/list", handler.Remind.List)
		admin.POST("/remind/delete", handler.Remind.Delete)
		admin.POST("/user/get", handler.User.Get)
		admin.POST("/user/post", handler.User.Post)
		admin.POST("/user/list", handler.User.List)
		admin.POST("/user/status", handler.User.Status)
	}

	// debug handler
	gee.DebugRoute(router)
	return router
}

var ProviderSet = wire.NewSet(NewRouter)
