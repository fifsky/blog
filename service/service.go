package service

import (
	"database/sql"
	"net/http"

	"app/config"
	"app/pkg/wechat"
	"app/store"
)

type Service struct {
	User    *User
	Article *Article
	Cate    *Cate
	Link    *Link
	Mood    *Mood
	Remind  *Remind
	Setting *Setting
	Weixin  *Weixin
}

func New(db *sql.DB, conf *config.Config, robot *wechat.Robot, httpClient *http.Client) *Service {
	s := store.New(db)
	return &Service{
		User:    NewUser(s, conf),
		Article: NewArticle(s, conf),
		Cate:    NewCate(s),
		Link:    NewLink(s),
		Mood:    NewMood(s),
		Remind:  NewRemind(s, robot, conf),
		Setting: NewSetting(s),
		Weixin:  NewWeixin(s, conf, httpClient),
	}
}
