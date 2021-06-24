package repo

import (
	"app/provider/model"
	"github.com/ilibs/gosql/v2"
)

type Remind struct {
	Base
}

func NewRemind(db *gosql.DB) *Remind {
	return &Remind{
		Base: Base{db: db},
	}
}

func (r *Remind) RemindGetList(start int, num int) ([]*model.Reminds, error) {
	var m = make([]*model.Reminds, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
