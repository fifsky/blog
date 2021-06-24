package repo

import (
	"app/provider/model"
	"github.com/ilibs/gosql/v2"
)

type User struct {
	Base
}

func NewUser(db *gosql.DB) *User {
	return &User{
		Base: Base{db: db},
	}
}

func (u *User) GetUser(uid int) (*model.Users, error) {
	user := &model.Users{}
	err := u.db.Model(user).Where("id = ?", uid).Get()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetList(start int, num int) ([]*model.Users, error) {
	var m = make([]*model.Users, 0)
	start = (start - 1) * num
	err := u.db.Model(&m).OrderBy("id desc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}
