package openapi

import (
	"net/http"

	"app/config"
	"app/service/feishu"
	"app/store"
)

type Service struct {
	User      *User
	Article   *Article
	Cate      *Cate
	Link      *Link
	Mood      *Mood
	Setting   *Setting
	Travel    *Travel
	MiniApp   *MiniApp
	Geo       *Geo
	Guestbook *Guestbook
	Comment   *Comment
}

func New(s *store.Store, conf *config.Config, httpClient *http.Client) *Service {
	sender := feishu.NewSender(conf.Feishu)
	return &Service{
		User:      NewUser(s, conf, sender, httpClient),
		Article:   NewArticle(s, conf),
		Cate:      NewCate(s),
		Link:      NewLink(s, conf, sender),
		Mood:      NewMood(s),
		Setting:   NewSetting(s),
		Travel:    NewTravel(s),
		MiniApp:   NewMiniApp(s, conf, httpClient),
		Geo:       NewGeo(s),
		Guestbook: NewGuestbook(s),
		Comment:   NewComment(s, conf, sender),
	}
}
