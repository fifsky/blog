package testdata

import (
	"database/sql"
	"fmt"
	"testing"

	"app/pkg/dbunit"
	"app/testutil"

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

// dump users testdata
func TestFixtures(t *testing.T) {
	if testing.Short() {
		return
	}

	// 使用 SQLite 数据库路径
	dbPath := "storage/blog.db"

	t.Run("users", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("users"), "select * from users limit 10")
		if err != nil {
			panic(err)
		}
	})

	t.Run("posts", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("posts"), "select * from posts where user_id = 1 limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("cates", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("cates"), "select * from cates where id = 1")
		if err != nil {
			panic(err)
		}
	})

	t.Run("comments", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("comments"), "select * from comments where post_id in(4,7)")
		if err != nil {
			panic(err)
		}
	})

	t.Run("links", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("links"), "select * from links limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("moods", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("moods"), "select * from moods limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("options", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("options"), "select * from options")
		if err != nil {
			panic(err)
		}
	})

	t.Run("reminds", func(t *testing.T) {
		_, err := dbunit.Dump(db(dbPath), testutil.Fixture("reminds"), "select * from reminds limit 2")
		if err != nil {
			panic(err)
		}
	})
}
