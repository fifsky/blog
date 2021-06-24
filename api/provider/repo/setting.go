package repo

import (
	"app/provider/model"
	"github.com/ilibs/gosql/v2"
)

type Setting struct {
	Base
}

func NewSetting(db *gosql.DB) *Setting {
	return &Setting{Base: Base{db: db}}
}

func (s *Setting) GetOptions() (map[string]string, error) {
	var options = make([]*model.Options, 0)
	err := s.db.Model(&options).All()
	if err != nil {
		return nil, err
	}

	options2 := make(map[string]string)
	for _, v := range options {
		options2[v.OptionKey] = v.OptionValue
	}
	return options2, nil
}
