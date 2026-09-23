package sqlext

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

const dbTag = "db"

var (
	// ErrTooManyColumns reports a query that returned several columns for a
	// destination of a primitive slice type, which can hold a single value.
	// For example, binding `select col1, col2 from mutable` to []string.
	ErrTooManyColumns = errors.New("too many columns returned for primitive slice")

	// ErrRowIntoSlice reports ScanRow being given a slice destination: it scans
	// a single row, so the destination has to be a pointer to a single value.
	ErrRowIntoSlice = errors.New("cannot scan a single row into a slice")

	// ErrNotPointer reports a destination that is not a pointer.
	ErrNotPointer = errors.New("not a pointer")

	// ErrNotSlicePointer reports a destination that is not a pointer to a slice.
	ErrNotSlicePointer = errors.New("not a pointer to a slice")

	// ColumnsMapper transforms struct field names into the database column
	// names. It is only used for fields without a db struct tag.
	// By default CamelCase field names are converted into snake_case:
	// E.g. FirstName -> first_name
	ColumnsMapper = func(name string) string { return snakeCase(name) }
)

// ScanRow scans a single row into a single variable. It requires that you use
// db.Query and not db.QueryRow, because QueryRow does not return column names.
// There is no performance impact in using one over the other. QueryRow only
// defers returning err until Scan is called, which is an unnecessary
// optimization for this library.
//
// The rows are not closed by ScanRow, the caller is responsible for closing them.
func ScanRow(v any, r *sql.Rows) error {
	vType := reflect.TypeOf(v)
	if k := vType.Kind(); k != reflect.Pointer {
		return fmt.Errorf("destination %q: %w", k.String(), ErrNotPointer)
	}

	vType = vType.Elem()
	vVal := reflect.ValueOf(v).Elem()
	if vType.Kind() == reflect.Slice {
		return ErrRowIntoSlice
	}

	sl := reflect.New(reflect.SliceOf(vType))
	err := ScanRows(sl.Interface(), r)
	if err != nil {
		return err
	}

	sl = sl.Elem()

	if sl.Len() == 0 {
		return sql.ErrNoRows
	}

	vVal.Set(sl.Index(0))

	return nil
}

// ScanRows scans sql rows into a slice (v). The rows are not closed by ScanRows,
// the caller is responsible for closing them.
func ScanRows(v any, r *sql.Rows) error {
	vType := reflect.TypeOf(v)
	if k := vType.Kind(); k != reflect.Pointer {
		return fmt.Errorf("destination %q: %w", k.String(), ErrNotPointer)
	}
	sliceType := vType.Elem()
	if reflect.Slice != sliceType.Kind() {
		return fmt.Errorf("destination %q: %w", sliceType.String(), ErrNotSlicePointer)
	}

	sliceVal := reflect.Indirect(reflect.ValueOf(v))
	itemType := sliceType.Elem()

	cols, err := r.Columns()
	if err != nil {
		return err
	}

	isPrimitive := itemType.Kind() != reflect.Struct

	for r.Next() {
		sliceItem := reflect.New(itemType).Elem()

		var pointers []any
		if isPrimitive {
			if len(cols) > 1 {
				return ErrTooManyColumns
			}
			pointers = []any{sliceItem.Addr().Interface()}
		} else {
			pointers = structPointers(sliceItem, cols)
		}

		if len(pointers) == 0 {
			return nil
		}

		err := r.Scan(pointers...)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, sliceItem))
	}
	return r.Err()
}

// Initialization the tags from struct. Fields are indexed by their column
// names: the db struct tag when it is set, otherwise the snake_case form of
// the field name. Fields tagged with `db:"-"` are ignored.
//
// Struct fields (and pointers to structs) that are not sql values themselves
// are recursed into: their fields are matched exactly like top level fields.
// A db tag on such a nested field is used as a column prefix, and a nil
// pointer is allocated in place. This makes self joins scannable:
//
//	var row struct {
//	    ID     int    `db:"id"`
//	    User   *User  `db:"u_"`   // matches u.id AS u_id, u.name AS u_name...
//	    Editor *User  `db:"le_"`  // matches le.id AS le_id...
//	}
func initFieldTag(sliceItem reflect.Value, fieldTagMap *map[string]reflect.Value) {
	initFieldTags(sliceItem, fieldTagMap, "")
}

func initFieldTags(sliceItem reflect.Value, fieldTagMap *map[string]reflect.Value, prefix string) {
	typ := sliceItem.Type()
	for i := 0; i < typ.NumField(); i++ {
		valField := sliceItem.Field(i)
		if !valField.IsValid() || !valField.CanSet() {
			continue
		}

		field := typ.Field(i)
		if ignored(field) {
			continue
		}

		if nested, ok := nestedStruct(valField); ok {
			// found a nested or embedded struct, the db tag (if any)
			// prefixes the columns of its fields
			initFieldTags(nested, fieldTagMap, prefix+fieldPrefix(field))
			continue
		}

		(*fieldTagMap)[prefix+columnName(field)] = valField
	}
}

// nestedStruct returns the struct behind a struct or struct pointer field
// when sqlext should recurse into it, allocating nil pointers in place.
// Structs that are valid sql values themselves (time.Time, driver.Valuer
// implementations) are scanned as single column and not recursed into.
func nestedStruct(valField reflect.Value) (reflect.Value, bool) {
	v := valField
	if v.Kind() == reflect.Pointer {
		if v.Type().Elem().Kind() != reflect.Struct || isValidSqlValue(v) {
			return reflect.Value{}, false
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return v.Elem(), true
	}

	if v.Kind() != reflect.Struct || isValidSqlValue(v) {
		return reflect.Value{}, false
	}

	return v, true
}

// fieldPrefix returns the column prefix a nested struct contributes to its
// fields: its db tag when one is set, no prefix otherwise.
func fieldPrefix(field reflect.StructField) string {
	if tag, ok := field.Tag.Lookup(dbTag); ok {
		return tag
	}
	return ""
}

func structPointers(sliceItem reflect.Value, cols []string) []any {
	pointers := make([]any, 0, len(cols))
	fieldTag := make(map[string]reflect.Value, sliceItem.NumField())
	initFieldTag(sliceItem, &fieldTag)

	for _, colName := range cols {
		fieldVal, ok := fieldTag[colName]
		if !ok || !fieldVal.IsValid() || !fieldVal.CanSet() {
			// have to add if we found a column because Scan() requires
			// len(cols) arguments or it will error. This way we can scan to
			// a useless pointer
			var nothing any
			pointers = append(pointers, &nothing)
			continue
		}

		pointers = append(pointers, fieldVal.Addr().Interface())
	}
	return pointers
}

// ignored reports whether a field is excluded from columns and scanning
// with the `db:"-"` struct tag.
func ignored(field reflect.StructField) bool {
	tag, ok := field.Tag.Lookup(dbTag)
	return ok && tag == "-"
}

// columnName returns the database column name of a field. The db struct tag
// takes precedence, otherwise the snake_case form of the field name is used.
func columnName(field reflect.StructField) string {
	if tag, ok := field.Tag.Lookup(dbTag); ok && tag != "" {
		return tag
	}
	return ColumnsMapper(field.Name)
}

func isValidSqlValue(v reflect.Value) bool {
	// This method covers two cases in which we know the Value can be converted to sql:
	// 1. It returns true for sql.driver's type check for types like time.Time
	// 2. It implements the driver.Valuer interface allowing conversion directly
	//    into sql statements
	if v.Kind() == reflect.Pointer {
		ptrVal := reflect.New(v.Type().Elem())
		return isValidSqlValue(ptrVal.Elem())
	}

	if driver.IsValue(v.Interface()) {
		return true
	}

	valuerType := reflect.TypeFor[driver.Valuer]()
	return v.Type().Implements(valuerType)
}

// snakeCase converts a CamelCase name into its snake_case representation.
// E.g. FirstName -> first_name, UserID -> user_id, HTTPServer -> http_server
func snakeCase(name string) string {
	var b strings.Builder
	b.Grow(len(name) + 4)

	var prev rune
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(prev) || unicode.IsDigit(prev) ||
				unicode.IsUpper(prev) && startsLower(name[i+utf8.RuneLen(r):])) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}

		prev = r
	}

	return b.String()
}

// startsLower reports whether s starts with a lower case letter.
func startsLower(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsLower(r)
}
