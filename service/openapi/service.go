package openapi

import (
	"net/http"

	"app/config"
	"app/store"
)

type Service struct {
	User      *User
	Article   *Article
	Cate      *Cate
	Link      *Link
	Mood      *Mood
	Remind    *Remind
	Setting   *Setting
	Travel    *Travel
	MiniApp   *MiniApp
	Geo       *Geo
	Guestbook *Guestbook
}

func New(s *store.Store, conf *config.Config, httpClient *http.Client) *Service {
	return &Service{
		User:      NewUser(s, conf),
		Article:   NewArticle(s, conf),
		Cate:      NewCate(s),
		Link:      NewLink(s),
		Mood:      NewMood(s),
		Remind:    NewRemind(s, conf),
		Setting:   NewSetting(s),
		Travel:    NewTravel(s),
		MiniApp:   NewMiniApp(s, conf, httpClient),
		Geo:       NewGeo(s),
		Guestbook: NewGuestbook(s),
	}
}
