package fixtures

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"
)

func TestFixtureFile(t *testing.T) {
	f := &fixtureFile{fileName: "posts.yml"}
	file := f.fileNameWithoutExtension()
	if file != "posts" {
		t.Errorf("Should be 'posts', but returned %s", file)
	}
}

func TestRequiredOptions(t *testing.T) {
	t.Run("DatabaseIsRequired", func(t *testing.T) {
		_, err := New()
		if err != errDatabaseIsRequired {
			t.Error("should return an error if database if not given")
		}
	})
}

const sqliteSchema = `CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_name TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  real_name TEXT NOT NULL DEFAULT '',
  password TEXT NOT NULL DEFAULT '',
  avatar TEXT NOT NULL DEFAULT '',
  status INTEGER NOT NULL DEFAULT 1,
  about TEXT NOT NULL DEFAULT '',
  role TEXT NOT NULL DEFAULT 'user',
  organization TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS un_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS un_user_name ON users(user_name);`

func TestLoader_Load(t *testing.T) {
	// 使用命名内存 SQLite 数据库进行测试
	db, err := sql.Open("sqlite", "file:test_loader?mode=memory&cache=shared")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(sqliteSchema)
	require.NoError(t, err)

	options := []func(*Loader) error{
		Database(db),
		Files("../testdata/fixtures/users.yml"),
		SkipTestDatabaseCheck(), // 内存数据库跳过测试库检查
	}

	f, err := New(options...)
	require.NoError(t, err)

	err = f.Load()
	require.NoError(t, err)

	row := db.QueryRow("select email from users where id = 1")
	var content string
	err = row.Scan(&content)
	require.NoError(t, err)
	require.Equal(t, `test@test.cn`, content)
}

func Test_sqliteHelper_quoteKeyword(t *testing.T) {
	h := &sqliteHelper{}
	k := h.quoteKeyword("status")
	require.Equal(t, "`status`", k)
}
