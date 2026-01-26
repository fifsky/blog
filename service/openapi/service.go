package openapi

import (
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
	Travel  *Travel
	MiniApp *MiniApp
	Geo     *Geo
}

func New(s *store.Store, conf *config.Config, robot *wechat.Robot, httpClient *http.Client) *Service {
	return &Service{
		User:    NewUser(s, conf),
		Article: NewArticle(s, conf),
		Cate:    NewCate(s),
		Link:    NewLink(s),
		Mood:    NewMood(s),
		Remind:  NewRemind(s, robot, conf),
		Setting: NewSetting(s),
		Weixin:  NewWeixin(s, conf, httpClient),
		Travel:  NewTravel(s),
		MiniApp: NewMiniApp(s, conf, httpClient),
		Geo:     NewGeo(s),
	}
}
