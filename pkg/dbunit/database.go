package dbunit

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite"
)

var (
	createTableRegex       = regexp.MustCompile(`(?isU)CREATE TABLE\s+.*;`)
	id               int32 = 0
	testDir                = filepath.Join(os.TempDir(), "blog_dbunit")
)

func init() {
	// 确保测试目录存在
	_ = os.MkdirAll(testDir, 0755)
}

// SetDatabase 设置测试数据库 DSN（SQLite 兼容，仅为兼容旧 API 保留）
func SetDatabase(dsn string) {
	// no-op: SQLite 使用文件路径，不需要 MySQL 的 server DSN
}

type database struct {
	Name   string // 数据库文件名（含 test_ 前缀）
	source string // 完整文件路径
	db     *sql.DB
}

func newDatabase(schema string) *database {
	atomic.AddInt32(&id, 1)
	name := fmt.Sprintf("test_%d_%d.db", time.Now().UnixNano(), id)
	return newDatabaseWithName(name, schema)
}

func newDatabaseWithName(name string, schema string) *database {
	dbPath := filepath.Join(testDir, name)
	d := &database{Name: name, source: dbPath}

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

// DSN 返回 SQLite 连接字符串（含 test_ 前缀，满足 EnsureTestDatabase 检查）
func (d *database) DSN() string {
	return fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)", d.source)
}

func (d *database) connection() error {
	db, err := sql.Open("sqlite", d.DSN())
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

// Drop 删除 SQLite 数据库文件
func (d *database) Drop() error {
	if d.db != nil {
		_ = d.db.Close()
	}
	// 删除主数据库文件及 WAL/SHM 文件
	for _, suffix := range []string{"", "-wal", "-shm"} {
		_ = os.Remove(d.source + suffix)
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

	db, err := sql.Open("sqlite", d.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	for _, query := range queries {
		if len(query) > 0 {
			if _, err := db.Exec(query); err != nil {
				return err
			}
		}
	}
	return nil
}
