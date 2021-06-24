package repo

import (
	"app/provider/model"
	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type Link struct {
	Base
}

func NewLink(db *gosql.DB) *Link {
	return &Link{Base: Base{
		db: db,
	}}
}

func (l *Link) GetAllLinks() []*model.Links {
	links := make([]*model.Links, 0)
	err := l.db.Model(&links).All()
	if err != nil {
		logger.Error(err)
		return nil
	}
	return links
}
