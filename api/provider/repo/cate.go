package repo

import (
	"app/provider/model"
	"github.com/ilibs/gosql/v2"
)

type Cate struct {
	Base
}

func NewCate(db *gosql.DB) *Cate {
	return &Cate{Base: Base{db: db}}
}

type CateArtivleCount struct {
	model.Cates
	Num int `json:"num" db:"num"`
}

func (a *Cate) GetAllCates() ([]*CateArtivleCount, error) {
	var cates = make([]*CateArtivleCount, 0)

	err := a.db.Select(&cates, "select c.*,ifnull(p.num,0) num from cates c left join (select count(*) num ,cate_id from posts where status = 1 and type = 1 group by cate_id) p on c.id = p.cate_id")
	if err != nil {
		return nil, err
	}
	return cates, nil
}
