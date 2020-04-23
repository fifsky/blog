package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/goapt/gee"

	"app/handler"
	"app/router/middleware"
)

func Route(router *gee.Engine) {

	// if CORS the remove annotation
	router.Use(gee.Wrap(cors.New(cors.Config{
		AllowOrigins:     []string{"http://fifsky.com", "http://www.fifsky.com", "https://fifsky.com", "https://www.fifsky.com"},
		AllowHeaders:     []string{"*"},
		MaxAge:           24 * time.Hour,
		AllowCredentials: false,
	})))

	// 中间件
	router.Use(gee.Wrap(middleware.Ginrus()))

	router.NoRoute(handler.Handle404)

	router.POST("/api/dingmsg", handler.DingMsg)
	router.POST("/api/login", handler.Login)
	router.POST("/api/mood/list", handler.MoodList)
	router.POST("/api/cate/all", handler.CateAll)
	router.POST("/api/article/archive", handler.ArchiveArticle)
	router.POST("/api/article/prevnext", handler.PrevNextArticle)
	router.POST("/api/article/list", handler.ListArticle)
	router.POST("/api/article/detail", handler.DetailArticle)
	router.POST("/api/link/all", handler.LinkAll)
	router.POST("/api/comment/new", handler.NewComment)
	router.POST("/api/comment/list", handler.CommentList)
	router.POST("/api/comment/post", handler.CommentPost)
	router.GET("/api/avatar", handler.Avatar)
	router.GET("/api/remind/change", middleware.RemindAuth, handler.RemindChange)
	router.GET("/api/remind/delay", middleware.RemindAuth, handler.RemindDelay)
	router.POST("/api/setting", handler.AdminSetting)
	router.GET("/feed.xml", handler.FeedGet)

	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthLogin)
	{
		admin.POST("/loginuser", handler.AdminLoginUser)
		admin.POST("/article/post", handler.AdminArticlePost)
		admin.POST("/article/delete", handler.AdminArticleDelete)
		admin.POST("/upload", handler.AdminUploadPost)
		admin.POST("/setting/post", handler.AdminSettingPost)
		admin.POST("/comment/list", handler.AdminCommentList)
		admin.POST("/comment/delete", handler.AdminCommentDelete)
		admin.POST("/mood/post", handler.AdminMoodPost)
		admin.POST("/mood/delete", handler.AdminMoodDelete)
		admin.POST("/cate/post", handler.AdminCatePost)
		admin.POST("/cate/list", handler.AdminCateList)
		admin.POST("/cate/delete", handler.AdminCateDelete)
		admin.POST("/link/post", handler.AdminLinkPost)
		admin.POST("/link/list", handler.AdminLinkList)
		admin.POST("/link/delete", handler.AdminLinkDelete)
		admin.POST("/remind/post", handler.AdminRemindPost)
		admin.POST("/remind/list", handler.AdminRemindList)
		admin.POST("/remind/delete", handler.AdminRemindDelete)
		admin.POST("/user/get", handler.AdminUserGet)
		admin.POST("/user/post", handler.AdminUserPost)
		admin.POST("/user/list", handler.AdminUserList)
		admin.POST("/user/status", handler.AdminUserStatus)
	}

	// debug handler
	gee.DebugRoute(router.Engine)
}
