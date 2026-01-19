// Package sqlutil provides utility functions for database operations with generics support.
package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// DB is the interface for database operations, compatible with *sql.DB and *sql.Tx.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Exec executes a query without returning any rows.
// It returns the number of rows affected.
func Exec(ctx context.Context, db DB, query string, args ...any) (int64, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("exec: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	return affected, nil
}

// ExecResult executes a query and returns the full sql.Result.
func ExecResult(ctx context.Context, db DB, query string, args ...any) (sql.Result, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}
	return result, nil
}

// Query executes a query and returns a slice of T.
// T must be a struct type with db tags on its fields.
func Query[T any](ctx context.Context, db DB, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("columns: %w", err)
	}

	var result []T
	for rows.Next() {
		var item T
		if err := scanRow(&item, rows, columns); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}

// QueryOne executes a query and returns a single result.
// Returns sql.ErrNoRows if no rows found.
func QueryOne[T any](ctx context.Context, db DB, query string, args ...any) (T, error) {
	var zero T
	results, err := Query[T](ctx, db, query, args...)
	if err != nil {
		return zero, err
	}
	if len(results) == 0 {
		return zero, sql.ErrNoRows
	}
	return results[0], nil
}

// scanRow scans a single row into the given struct pointer.
// It handles NULL values for non-pointer types by returning their zero value.
func scanRow(dest any, rows *sql.Rows, columns []string) error {
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	// Build a map from db tag to field index
	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		if tag == "-" {
			continue
		}
		fieldMap[tag] = i
	}

	// Create scan destinations with nullable wrappers
	scanDest := make([]any, len(columns))
	fieldIndices := make([]int, len(columns))  // Track which field each column maps to
	needsConvert := make([]bool, len(columns)) // Track if we need to convert from nullable

	for i, col := range columns {
		fieldIndices[i] = -1
		if fieldIdx, ok := fieldMap[col]; ok {
			fieldIndices[i] = fieldIdx
			field := v.Field(fieldIdx)
			fieldType := field.Type()

			// Check if the field is already a sql.Null* type or pointer
			if isNullableType(fieldType) || fieldType.Kind() == reflect.Ptr {
				// Use the field directly
				scanDest[i] = field.Addr().Interface()
			} else {
				// Use a nullable wrapper to handle NULL values
				scanDest[i] = createNullableWrapper(fieldType)
				needsConvert[i] = true
			}
		} else {
			// Column not found in struct, use a throwaway destination
			var discard any
			scanDest[i] = &discard
		}
	}

	if err := rows.Scan(scanDest...); err != nil {
		return err
	}

	// Convert nullable wrappers back to struct fields
	for i, idx := range fieldIndices {
		if idx >= 0 && needsConvert[i] {
			setFieldFromNullable(v.Field(idx), scanDest[i])
		}
	}

	return nil
}

// isNullableType checks if a type is a sql.Null* type
func isNullableType(t reflect.Type) bool {
	switch t {
	case reflect.TypeOf(sql.NullString{}),
		reflect.TypeOf(sql.NullInt64{}),
		reflect.TypeOf(sql.NullInt32{}),
		reflect.TypeOf(sql.NullInt16{}),
		reflect.TypeOf(sql.NullFloat64{}),
		reflect.TypeOf(sql.NullBool{}),
		reflect.TypeOf(sql.NullTime{}),
		reflect.TypeOf(sql.NullByte{}):
		return true
	}
	return false
}

// createNullableWrapper creates a nullable wrapper for the given type
func createNullableWrapper(t reflect.Type) any {
	switch t.Kind() {
	case reflect.String:
		return new(sql.NullString)
	case reflect.Int, reflect.Int64:
		return new(sql.NullInt64)
	case reflect.Int32:
		return new(sql.NullInt32)
	case reflect.Int16:
		return new(sql.NullInt16)
	case reflect.Int8:
		return new(sql.NullInt16) // Use NullInt16 for int8
	case reflect.Float64:
		return new(sql.NullFloat64)
	case reflect.Float32:
		return new(sql.NullFloat64) // Use NullFloat64 for float32
	case reflect.Bool:
		return new(sql.NullBool)
	case reflect.Uint8:
		return new(sql.NullByte)
	default:
		// For complex types (time.Time, []byte, etc.), use interface{}
		return new(any)
	}
}

// setFieldFromNullable sets the struct field value from a nullable wrapper
func setFieldFromNullable(field reflect.Value, wrapper any) {
	switch w := wrapper.(type) {
	case *sql.NullString:
		if w.Valid {
			field.SetString(w.String)
		}
	case *sql.NullInt64:
		if w.Valid {
			field.SetInt(w.Int64)
		}
	case *sql.NullInt32:
		if w.Valid {
			field.SetInt(int64(w.Int32))
		}
	case *sql.NullInt16:
		if w.Valid {
			field.SetInt(int64(w.Int16))
		}
	case *sql.NullFloat64:
		if w.Valid {
			field.SetFloat(w.Float64)
		}
	case *sql.NullBool:
		if w.Valid {
			field.SetBool(w.Bool)
		}
	case *sql.NullByte:
		if w.Valid {
			field.SetUint(uint64(w.Byte))
		}
	case *any:
		if *w != nil {
			val := reflect.ValueOf(*w)
			if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			}
		}
	}
}

