package dbunit

import (
	"database/sql"
	"fmt"
	"io"

	"app/pkg/dbunit/fixtures"

	_ "modernc.org/sqlite"
)

type Testing struct {
	tdb *database
	db  *sql.DB
}

func NewTest(schema io.Reader) *Testing {
	tdb := newDatabase(schema)

	// 打开 SQLite 测试数据库连接
	db, err := sql.Open("sqlite", tdb.DSN())
	if err != nil {
		panic("test sqlite open fail " + err.Error())
	}

	return &Testing{
		tdb: tdb,
		db:  db,
	}
}

func (d *Testing) DB() *sql.DB {
	return d.db
}

// Drop 关闭测试数据库连接，内存数据库自动释放
func (d *Testing) Drop() {
	// 先关闭数据库连接
	if d.db != nil {
		_ = d.db.Close()
	}
	err := d.tdb.Drop()
	if err != nil {
		panic("drop database error " + err.Error())
	}
}

// Load 加载 YAML fixture 数据到测试数据库
// data 为文件名到内容的映射，文件名需包含 .yml 扩展名（用于推断表名）
func (d *Testing) Load(data map[string][]byte) {
	if len(data) == 0 {
		return
	}

	options := make([]func(*fixtures.Loader) error, 0)
	options = append(options, fixtures.Database(d.db))
	options = append(options, fixtures.SkipTestDatabaseCheck()) // 内存数据库跳过测试库检查
	options = append(options, fixtures.Content(data))

	f, err := fixtures.New(options...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("db Load fixtures:%d files\n", len(data))

	if err := f.Load(); err != nil {
		panic(err)
	}
}
