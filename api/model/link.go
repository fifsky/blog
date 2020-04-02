package model

import (
	"time"

	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type Links struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Url       string    `form:"url" json:"url" db:"url"`
	Desc      string    `form:"desc" json:"desc" db:"desc"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (l *Links) TableName() string {
	return "links"
}

func (l *Links) PK() string {
	return "id"
}

func (p *Links) AfterChange() {
	Cache.Delete("all-links")
}

func GetAllLinks() []*Links {
	if v, ok := Cache.Get("all-links"); ok {
		return v.([]*Links)
	}
	links := make([]*Links, 0)
	err := gosql.Model(&links).All()
	if err != nil {
		logger.Error(err)
		return nil
	}

	Cache.Set("all-links", links, 1*time.Hour)
	return links
}
