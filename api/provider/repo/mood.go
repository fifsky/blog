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

	err := gosql.Model(&moods).Limit(num).Offset(start).OrderBy("id desc").All()

	if err != nil {
		return nil, err
	}

	return moods, err
}
