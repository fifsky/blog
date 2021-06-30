package repo

import (
	"app/provider/model"
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

func (l *Link) GetAllLinks() ([]*model.Links, error) {
	links := make([]*model.Links, 0)
	err := l.db.Model(&links).All()
	return links, err
}
