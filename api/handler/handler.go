package handler

import (
	"database/sql"

	"app/config"
	"app/pkg/wechat"
	"app/store"
)

type Handler struct {
	Article *Article
	Cate    *Cate
	Link    *Link
	User    *User
	Mood    *Mood
	Remind  *Remind
	Setting *Setting
}

func New(db *sql.DB, conf *config.Config, robot *wechat.Robot) *Handler {
	s := store.New(db)
	return &Handler{
		Article: NewArticle(s, conf),
		Cate:    NewCate(s),
		Link:    NewLink(s),
		User:    NewUser(s, conf),
		Mood:    NewMood(s),
		Remind:  NewRemind(s, robot),
		Setting: NewSetting(s),
	}
}
