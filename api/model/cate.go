package model

import (
	"time"

	"github.com/goapt/logger"
	"github.com/ilibs/gosql/v2"
)

type Cates struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Name      string    `form:"name" json:"name" db:"name"`
	Desc      string    `form:"desc" json:"desc" db:"desc"`
	Domain    string    `form:"domain" json:"domain" db:"domain"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
	UpdatedAt time.Time `form:"-" json:"updated_at" time_format:"2006-01-02 15:04:05" db:"updated_at"`
}

func (c *Cates) TableName() string {
	return "cates"
}

func (c *Cates) PK() string {
	return "id"
}

func (c *Cates) AfterChange() {
	Cache.Delete("all-cates")
}

type CateArtivleCount struct {
	Cates
	Num int `json:"num" db:"num"`
}

func GetAllCates() []*CateArtivleCount {
	if v, ok := Cache.Get("all-cates"); ok {
		return v.([]*CateArtivleCount)
	}
	var cates = make([]*CateArtivleCount, 0)

	err := gosql.Select(&cates, "select c.*,ifnull(p.num,0) num from cates c left join (select count(*) num ,cate_id from posts where status = 1 and type = 1 group by cate_id) p on c.id = p.cate_id")
	if err != nil {
		logger.Error(err)
		return nil
	}

	Cache.Set("all-cates", cates, 1*time.Hour)
	return cates
}
