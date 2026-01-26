package admin

import (
	"app/config"
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
	AI      *AI
	Photo   *Photo
	OSS     *OSS
	Region  *Region
}

func New(s *store.Store, conf *config.Config) *Service {
	return &Service{
		User:    NewUser(s, conf),
		Article: NewArticle(s, conf),
		Cate:    NewCate(s),
		Link:    NewLink(s),
		Mood:    NewMood(s),
		Remind:  NewRemind(s),
		Setting: NewSetting(s),
		AI:      NewAI(conf, s),
		Photo:   NewPhoto(s),
		OSS:     NewOSS(conf),
		Region:  NewRegion(s),
	}
}
