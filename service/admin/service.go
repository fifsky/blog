package admin

import (
	"app/config"
	"app/pkg/agent"

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
	AI        *AI
	OSS       *OSS
	ClawBot   *ClawBot
	Guestbook *Guestbook
	Footprint *Footprint
	Comment   *Comment
}

func New(s *store.Store, conf *config.Config, aiAgent *agent.Agent) *Service {
	return &Service{
		User:      NewUser(s, conf),
		Article:   NewArticle(s, conf),
		Cate:      NewCate(s),
		Link:      NewLink(s),
		Mood:      NewMood(s),
		Remind:    NewRemind(s),
		Setting:   NewSetting(s),
		AI:        NewAI(aiAgent, s, conf),
		OSS:       NewOSS(conf),
		ClawBot:   NewClawBot(s, conf, aiAgent),
		Guestbook: NewGuestbook(s),
		Footprint: NewFootprint(s),
		Comment:   NewComment(s),
	}
}
