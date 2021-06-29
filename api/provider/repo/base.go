package repo

import (
	"github.com/ilibs/gosql/v2"
)

type BuilderFunc func(b *gosql.Builder)

type Base struct {
	db *gosql.DB
}

func (b *Base) Create(m gosql.IModel) (int64, error) {
	return b.db.Model(m).Create()
}

func (b *Base) Update(m gosql.IModel) (int64, error) {
	return b.db.Model(m).Update()
}

func (b *Base) Find(m gosql.IModel) error {
	return b.db.Model(m).Get()
}

func (b *Base) FindAll(m interface{}, fn ...BuilderFunc) error {
	builder := b.db.Model(m)

	if len(fn) > 0 {
		for _, f := range fn {
			f(builder)
		}
	}

	return builder.All()
}
