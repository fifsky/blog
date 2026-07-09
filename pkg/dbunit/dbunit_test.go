package dbunit

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestRun(t *testing.T) {
	t.Run("with fixtures", func(t *testing.T) {
		Run(t, testSchemaReader(), func(t *testing.T, db *sql.DB) {
			row := db.QueryRow("select email from users where id = 1")
			var email string
			if err := row.Scan(&email); err != nil {
				t.Fatal(err)
			}

			if email != "test@test.cn" {
				t.Fatalf("user mismatch want %s,but get %s", "test@test.cn", email)
			}
		}, testFixturesMap("users"))
	})

	t.Run("select fixtures", func(t *testing.T) {
		Run(t, testSchemaReader(), func(t *testing.T, db *sql.DB) {
			row := db.QueryRow("select email from users where id = 1")
			var email string
			if err := row.Scan(&email); err != sql.ErrNoRows {
				t.Fatal(err)
			}
		}, testFixturesMap("members", "documents"))
	})

	t.Run("custom fixtures", func(t *testing.T) {
		Run(t, testSchemaReader(), func(t *testing.T, db *sql.DB) {
			var ct int
			err := db.QueryRow("select count(1) from custom").Scan(&ct)

			if err != nil {
				t.Fatal(err)
			}

			if ct == 0 {
				t.Fatalf("user mismatch want %s,but get %d", " > 0", ct)
			}
		}, testCustomFixtures())
	})
}

func TestNew(t *testing.T) {
	New(t, func(d *DBUnit) {
		db := d.NewDatabase(testSchemaReader(), testFixturesMap("users"))
		_ = d.NewDatabase(testSchemaReader(), nil)
		row := db.QueryRow("select email from users where id = 1")
		var email string
		if err := row.Scan(&email); err != nil {
			t.Fatal(err)
		}

		if email != "test@test.cn" {
			t.Fatalf("user mismatch want %s,but get %s", "test@test.cn", email)
		}
	})
}

func TestLoad(t *testing.T) {
	test := NewTest(testSchemaReader())
	t.Cleanup(func() {
		test.Drop()
	})

	test.Load(testCustomFixtures())
}
