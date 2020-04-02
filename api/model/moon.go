package model

import (
	"time"

	"github.com/ilibs/gosql/v2"
)

type Moods struct {
	Id        int       `form:"id" json:"id" db:"id"`
	Content   string    `form:"content" json:"content" db:"content"`
	UserId    int       `form:"user_id" json:"user_id" db:"user_id"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (m *Moods) TableName() string {
	return "moods"
}

func (m *Moods) PK() string {
	return "id"
}

func (m *Moods) AfterChange() {
	Cache.Delete("remind-first")
}

type UserMoods struct {
	Moods
	User *Users `json:"user" db:"nick_name" relation:"user_id,id"`
}

func MoodGetList(start int, num int) ([]*UserMoods, error) {

	var moods = make([]*UserMoods, 0)
	start = (start - 1) * num

	err := gosql.Model(&moods).Limit(num).Offset(start).OrderBy("id desc").All()

	if err != nil {
		return nil, err
	}

	return moods, err
}
