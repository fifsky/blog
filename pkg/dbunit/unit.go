package dbunit

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"app/pkg/dbunit/fixtures"

	_ "modernc.org/sqlite"
)

type Testing struct {
	tdb *database
	db  *sql.DB
}

func NewTest(schema string) *Testing {
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

// Load 加载 YAML fixture 文件到测试数据库
func (d *Testing) Load(files ...string) {
	options := make([]func(*fixtures.Loader) error, 0)
	options = append(options, fixtures.Database(d.db))
	options = append(options, fixtures.SkipTestDatabaseCheck()) // 内存数据库跳过测试库检查

	fs := make([]string, 0)
	for _, file := range files {
		if isDir(file) {
			options = append(options, fixtures.Directory(file))
		} else {
			fs = append(fs, file)
		}
	}

	if len(fs) > 0 {
		options = append(options, fixtures.Files(fs...))
	}

	f, err := fixtures.New(options...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("db Load fixtures:%s\n", strings.Join(files, ","))

	if err := f.Load(); err != nil {
		panic(err)
	}
}

// isDir 判断路径是否为目录
func isDir(path string) bool {
	fio, err := os.Lstat(path)
	if nil != err {
		return false
	}

	return fio.IsDir()
}
