package dbunit

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
)

func TestNewDatabase(t *testing.T) {
	tdb := newDatabase("./testdata/schema.sql")
	t.Cleanup(func() {
		_ = tdb.Drop()
	})

	db, err := sql.Open("sqlite", tdb.DSN())
	assert.NoError(t, err)
	defer db.Close()
}
