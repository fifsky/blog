package sqlext

import (
	"context"
	"database/sql"
)

// Queryer is the only method QueryRow and Query need to run a query. Both
// *sql.DB and *sql.Tx satisfy it.
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// QueryRow runs the query and scans the first row into T, returning
// sql.ErrNoRows when the query returns no row. The zero value of T is returned
// along with any error, so callers can ignore the value when err != nil.
//
// The rows are closed by QueryRow. Columns are matched with the struct fields
// by their db tag, or the snake_case form of the field name, see README.
func QueryRow[T any](ctx context.Context, db Queryer, query string, args ...any) (T, error) {
	var dest T

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return dest, err
	}
	defer rows.Close()

	if err := ScanRow(&dest, rows); err != nil {
		return dest, err
	}
	return dest, nil
}

// Query runs the query and scans every row into a []T. An empty result is an
// empty slice rather than a nil one, so it marshals to [] instead of null.
//
// The rows are closed by Query. Columns are matched with the struct fields by
// their db tag, or the snake_case form of the field name, see README.
func Query[T any](ctx context.Context, db Queryer, query string, args ...any) ([]T, error) {
	dest := make([]T, 0)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := ScanRows(&dest, rows); err != nil {
		return nil, err
	}
	return dest, nil
}
