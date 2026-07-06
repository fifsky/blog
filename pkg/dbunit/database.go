package dbunit

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite"
)

var (
	createTableRegex       = regexp.MustCompile(`(?isU)CREATE TABLE\s+.*;`)
	id               int32 = 0
)

type database struct {
	name string // 内存数据库名称（含 test_ 前缀）
	db   *sql.DB
}

func newDatabase(schema string) *database {
	atomic.AddInt32(&id, 1)
	name := fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), id)
	return newDatabaseWithName(name, schema)
}

func newDatabaseWithName(name string, schema string) *database {
	d := &database{name: name}

	err := d.connection()
	if err != nil {
		panic("test sqlite connection fail," + err.Error())
	}

	err = d.Import(schema)
	if err != nil {
		panic(err)
	}
	return d
}

// DSN 返回内存 SQLite 连接字符串
func (d *database) DSN() string {
	return fmt.Sprintf("file:%s?mode=memory&cache=shared&_pragma=foreign_keys(ON)", d.name)
}

func (d *database) connection() error {
	db, err := sql.Open("sqlite", d.DSN())
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

// Drop 关闭数据库连接，内存数据库在所有连接关闭后自动释放
func (d *database) Drop() error {
	if d.db != nil {
		_ = d.db.Close()
	}
	return nil
}

// Import 导入 schema 文件，执行所有 CREATE TABLE 语句
func (d *database) Import(schema string) error {
	if !isExists(schema) {
		return fmt.Errorf("sql file not found:%s", schema)
	}

	content, err := os.ReadFile(schema)
	if err != nil {
		return err
	}

	// 提取所有 CREATE TABLE 语句
	queries := createTableRegex.FindAllString(string(content), -1)

	// 同时提取 CREATE INDEX、CREATE UNIQUE INDEX、CREATE TRIGGER 语句
	indexRe := regexp.MustCompile(`(?isU)CREATE\s+(?:UNIQUE\s+)?INDEX\s+.*;`)
	triggerRe := regexp.MustCompile(`(?isU)CREATE\s+TRIGGER\s+(?:IF\s+NOT\s+EXISTS\s+)?[\s\S]*?END;`)
	queries = append(queries, indexRe.FindAllString(string(content), -1)...)
	queries = append(queries, triggerRe.FindAllString(string(content), -1)...)

	for _, query := range queries {
		if len(query) > 0 {
			if _, err := d.db.Exec(query); err != nil {
				return err
			}
		}
	}
	return nil
}
