package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Database 数据库连接配置
type Database struct {
	Driver       string `yaml:"driver"`
	Dsn          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifetime  int    `yaml:"max_lifetime"`
}

// Connect 建立 SQLite 数据库连接，自动创建数据目录并启用 WAL 模式
func (d *Database) Connect() *sql.DB {
	// 从 DSN 中提取文件路径，自动创建目录
	dbPath := d.ExtractDBPath()
	if dbPath != "" && dbPath != ":memory:" {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("[db] failed to create database directory: %s\n", err)
		}
	}

	log.Printf("[db] connect sqlite: %s\n", dbPath)
	sess, err := sql.Open(d.Driver, d.Dsn)
	if err != nil {
		log.Fatalf("[db] failed to open: %s\n", err)
	}

	if err = sess.Ping(); err != nil {
		_ = sess.Close()
		log.Fatalf("[db] failed to connect: %s\n", err)
	}

	sess.SetMaxOpenConns(d.MaxOpenConns)
	sess.SetMaxIdleConns(d.MaxIdleConns)
	if d.MaxLifetime > 0 {
		sess.SetConnMaxLifetime(time.Duration(d.MaxLifetime) * time.Second)
	}
	return sess
}

// ExtractDBPath 从 SQLite DSN 中提取数据库文件路径
// DSN 格式: file:storage/blog.db?_pragma=journal_mode(WAL)&...
func (d *Database) ExtractDBPath() string {
	dsn := d.Dsn
	if after, ok := strings.CutPrefix(dsn, "file:"); ok {
		dsn = after
	}
	// 去掉查询参数
	if idx := strings.Index(dsn, "?"); idx != -1 {
		dsn = dsn[:idx]
	}
	return dsn
}
