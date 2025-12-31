package service

import (
	"database/sql"

	"app/config"
	"app/pkg/wechat"
	"app/store"
)

type Service struct {
	User *User
}

func New(db *sql.DB, conf *config.Config, robot *wechat.Robot) *Service {
	s := store.New(db)
	return &Service{
		User: NewUser(s, conf),
	}
}
