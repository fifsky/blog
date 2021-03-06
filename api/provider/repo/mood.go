package repo

import (
	"app/provider/model"
	"github.com/ilibs/gosql/v2"
)

type Mood struct {
	Base
}

func NewMood(db *gosql.DB) *Mood {
	return &Mood{Base: Base{
		db: db,
	}}
}

type UserMoods struct {
	model.Moods
	User *model.Users `json:"user" db:"nick_name" relation:"user_id,id"`
}

func (m *Mood) MoodGetList(start int, num int) ([]*UserMoods, error) {

	var moods = make([]*UserMoods, 0)
	start = (start - 1) * num

	err := m.db.Model(gosql.NewModelWrapper(map[string]*gosql.DB{
		"default": m.db,
	}, &moods)).Limit(num).Offset(start).OrderBy("id desc").All()

	return moods, err
}
