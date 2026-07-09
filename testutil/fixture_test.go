package testutil

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"app/pkg/dbunit"

	_ "modernc.org/sqlite"
)

// db 打开 SQLite 数据库连接
func db(dbPath string) *sql.DB {
	d, err := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)", dbPath))
	if err != nil {
		panic(err)
	}
	return d
}

// TestFixtures 从真实数据库导出测试数据到 fixture 文件
// 使用 -short 跳过此测试（需要真实数据库）
func TestFixtures(t *testing.T) {
	if testing.Short() {
		return
	}

	// 使用 SQLite 数据库路径
	dbPath := "storage/blog.db"

	t.Run("users", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "users.yml"), "select * from users limit 10")
		if err != nil {
			panic(err)
		}
	})

	t.Run("posts", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "posts.yml"), "select * from posts where user_id = 1 limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("cates", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "cates.yml"), "select * from cates where id = 1")
		if err != nil {
			panic(err)
		}
	})

	t.Run("comments", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "comments.yml"), "select * from comments where post_id in(4,7)")
		if err != nil {
			panic(err)
		}
	})

	t.Run("links", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "links.yml"), "select * from links limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("moods", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "moods.yml"), "select * from moods limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("options", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "options.yml"), "select * from options")
		if err != nil {
			panic(err)
		}
	})

	t.Run("reminds", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), filepath.Join("testdata", "fixtures", "reminds.yml"), "select * from reminds limit 2")
		if err != nil {
			panic(err)
		}
	})
}
