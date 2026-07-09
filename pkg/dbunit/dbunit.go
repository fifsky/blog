package dbunit

import (
	"database/sql"
	"io"
	"testing"
)

func Run(t *testing.T, schema io.Reader, f func(t *testing.T, db *sql.DB), fixtures map[string][]byte) {
	New(t, func(d *DBUnit) {
		db := d.NewDatabase(schema, fixtures)
		f(t, db)
	})
}

type DBUnit struct {
	tests []*Testing
}

// NewDatabase 创建内存 SQLite 数据库，导入 schema 和 fixture 数据
// schema 为 io.Reader（如嵌入的 schema.sql），fixtures 为文件名到内容的映射
func (d *DBUnit) NewDatabase(schema io.Reader, fixtures map[string][]byte) *sql.DB {
	test := NewTest(schema)
	test.Load(fixtures)
	d.tests = append(d.tests, test)
	return test.DB()
}

func (d *DBUnit) drop() {
	for _, test := range d.tests {
		test.Drop()
	}
}

func New(t *testing.T, f func(d *DBUnit)) {
	dt := &DBUnit{}
	t.Cleanup(func() {
		dt.drop()
	})

	f(dt)
}
