package main

import (
	"database/sql"
	"fmt"
	"os"

	"app/pkg/dbunit"

	_ "modernc.org/sqlite"
)

func main() {
	// 使用 SQLite 数据库进行 fixture 导出
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "storage/blog.db"
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)", dbPath))
	if err != nil {
		panic(err)
	}

	// 文档数据集
	data, err := dbunit.Dump(db, "testdata/fixtures/documents.yml", "select * from documents limit 3")
	if err != nil {
		panic(err)
	}

	// 用户数据集
	userIds := dbunit.Pluck(data, "user_id")
	_, err = dbunit.Dump(db, "testdata/fixtures/users.yml", "select * from users where id in(?)", userIds)
	if err != nil {
		panic(err)
	}

	// members
	docIds := dbunit.Pluck(data, "doc_id")
	_, err = dbunit.Dump(db, "testdata/fixtures/members.yml", "select * from members where doc_id in(?)", docIds)
	if err != nil {
		panic(err)
	}
}
