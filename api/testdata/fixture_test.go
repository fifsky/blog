package testdata

import (
	"fmt"
	"testing"

	"app/testutil"
	"github.com/goapt/dbunit"
	"github.com/ilibs/gosql/v2"
)

// the connect development environment database
func db(dbname string) *gosql.DB {
	d, err := gosql.Open("mysql", fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=%s", dbname, "Asia%2FShanghai"))
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

	t.Run("users", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("users"), "select * from users limit 10")
		if err != nil {
			panic(err)
		}
	})

	t.Run("posts", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("posts"), "select * from posts where user_id = 1 limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("cates", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("cates"), "select * from cates where id = 1")
		if err != nil {
			panic(err)
		}
	})

	t.Run("comments", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("comments"), "select * from comments where post_id in(4,7)")
		if err != nil {
			panic(err)
		}
	})

	t.Run("links", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("links"), "select * from links limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("moods", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("moods"), "select * from moods limit 2")
		if err != nil {
			panic(err)
		}
	})

	t.Run("options", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("options"), "select * from options")
		if err != nil {
			panic(err)
		}
	})

	t.Run("reminds", func(t *testing.T) {
		_, err := dbunit.Dump(db("blog"), testutil.Fixture("reminds"), "select * from reminds limit 2")
		if err != nil {
			panic(err)
		}
	})
}
