package sqlext_test

import (
	"database/sql"
	"testing"
	"time"

	"app/pkg/sqlext"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRowsMapsSnakeCaseColumnsToFieldNames(t *testing.T) {
	var item struct {
		FirstName string
	}

	expected := "Brett Jones"
	rows := queryRows(t, []string{"first_name"}, []any{expected})

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Equal(t, expected, item.FirstName)
}

func TestRowsMapsSnakeCaseAcronyms(t *testing.T) {
	var item struct {
		UserID int64
	}

	rows := queryRows(t, []string{"user_id"}, []any{int64(7)})

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Equal(t, int64(7), item.UserID)
}

func TestRowsUsesTagName(t *testing.T) {
	expected := "Brett Jones"
	rows := queryRows(t, []string{"first_and_last_name"}, []any{expected})

	var item struct {
		FirstAndLastName string `db:"first_and_last_name"`
	}

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Equal(t, expected, item.FirstAndLastName)
}

func TestRowsIgnoresUnsetableColumns(t *testing.T) {
	expected := "Brett Jones"
	rows := queryRows(t, []string{"first_and_last_name"}, []any{expected})

	var item struct {
		// private, unsetable
		firstAndLastName string `db:"first_and_last_name"`
	}

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.NotEqual(t, expected, item.firstAndLastName)
}

func TestErrorsWhenScanFails(t *testing.T) {
	rows := queryRows(t, []string{"age"}, []any{"not a number"})

	var item struct {
		Age int
	}

	err := sqlext.ScanRow(&item, rows)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "age")
	assert.Contains(t, err.Error(), "converting")
}

func TestRowsErrorsWhenNotGivenAPointer(t *testing.T) {
	// 类型检查先于读取 rows，用不到真实结果集
	err := sqlext.ScanRows("hello", nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, sqlext.ErrNotPointer)
	assert.Contains(t, err.Error(), "pointer")
}

func TestRowsErrorsWhenNotGivenAPointerToSlice(t *testing.T) {
	var item struct{}

	err := sqlext.ScanRows(&item, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, sqlext.ErrNotSlicePointer)
	assert.Contains(t, err.Error(), "slice")
}

func TestDoesNothingWhenNextIsFalse(t *testing.T) {
	rows := queryRows(t, []string{"Name"})

	var items []struct {
		Name string
		Age  int
	}
	err := sqlext.ScanRows(&items, rows)
	assert.NoError(t, err)
	assert.Nil(t, items)
}

func TestIgnoresColumnsThatDoNotHaveFields(t *testing.T) {
	rows := queryRows(t, []string{"first", "last", "age"},
		[]any{"Brett", "Jones", int8(40)},
		[]any{"Fred", "Jones", int8(50)},
	)

	var items []struct {
		First string
		Last  string
	}

	require.NoError(t, sqlext.ScanRows(&items, rows))
	require.Len(t, items, 2)
	assert.Equal(t, "Brett", items[0].First)
	assert.Equal(t, "Jones", items[0].Last)
	assert.Equal(t, "Fred", items[1].First)
	assert.Equal(t, "Jones", items[1].Last)
}

func TestIgnoresFieldsThatDoNotHaveColumns(t *testing.T) {
	rows := queryRows(t, []string{"first", "age"},
		[]any{"Brett", int8(40)},
		[]any{"Fred", int8(50)},
	)

	var items []struct {
		First string
		Last  string
		Age   int8
	}

	require.NoError(t, sqlext.ScanRows(&items, rows))
	require.Len(t, items, 2)
	assert.EqualValues(t, "Brett", items[0].First)
	assert.EqualValues(t, "", items[0].Last)
	assert.EqualValues(t, 40, items[0].Age)

	assert.EqualValues(t, "Fred", items[1].First)
	assert.EqualValues(t, "", items[1].Last)
	assert.EqualValues(t, 50, items[1].Age)
}

func TestRowScansToPrimitiveType(t *testing.T) {
	expected := "Bob"
	rows := queryRows(t, []string{"name"}, []any{expected})

	var name string
	assert.NoError(t, sqlext.ScanRow(&name, rows))
	assert.Equal(t, expected, name)
}

func TestScansPrimitiveSlices(t *testing.T) {
	t.Run("ints", func(t *testing.T) {
		rows := queryRows(t, []string{"a"}, []any{1}, []any{2}, []any{3})

		var scanned []int
		require.NoError(t, sqlext.ScanRows(&scanned, rows))
		assert.Equal(t, []int{1, 2, 3}, scanned)
	})

	t.Run("strings", func(t *testing.T) {
		rows := queryRows(t, []string{"a"}, []any{"brett"}, []any{"fred"}, []any{"geoff"})

		var scanned []string
		require.NoError(t, sqlext.ScanRows(&scanned, rows))
		assert.Equal(t, []string{"brett", "fred", "geoff"}, scanned)
	})

	t.Run("bools", func(t *testing.T) {
		rows := queryRows(t, []string{"a"}, []any{true}, []any{false})

		var scanned []bool
		require.NoError(t, sqlext.ScanRows(&scanned, rows))
		assert.Equal(t, []bool{true, false}, scanned)
	})

	t.Run("floats", func(t *testing.T) {
		rows := queryRows(t, []string{"a"}, []any{1.0}, []any{1.1}, []any{1.2})

		var scanned []float64
		require.NoError(t, sqlext.ScanRows(&scanned, rows))
		assert.Equal(t, []float64{1, 1.1, 1.2}, scanned)
	})
}

func TestErrorsWhenMoreThanOneColumnForPrimitiveSlice(t *testing.T) {
	rows := queryRows(t, []string{"fname", "lname"}, []any{"brett", "jones"})

	var fnames []string

	err := sqlext.ScanRows(&fnames, rows)
	assert.EqualValues(t, sqlext.ErrTooManyColumns, err)
}

func TestErrorsWhenScanRowToSlice(t *testing.T) {
	var persons []struct {
		ID int
	}

	err := sqlext.ScanRow(&persons, nil)
	assert.EqualValues(t, sqlext.ErrRowIntoSlice, err)
}

func TestRowReturnsErrNoRowsWhenQueryHasNoRows(t *testing.T) {
	rows := queryRows(t, []string{"First"})

	var item struct {
		First string
	}

	assert.EqualValues(t, sql.ErrNoRows, sqlext.ScanRow(&item, rows))
}

func TestRowErrorsWhenItemIsNotAPointer(t *testing.T) {
	var item struct {
		First string
	}

	err := sqlext.ScanRow(item, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, sqlext.ErrNotPointer)
	assert.Contains(t, err.Error(), "pointer")
}

func TestRowUsesDBTagOverFieldName(t *testing.T) {
	rows := queryRows(t, []string{"first", "last"}, []any{"Brett", "Jones"})

	var item struct {
		First string `db:"first"`
		Last  string `db:"-"`
	}

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Equal(t, "Brett", item.First)
	assert.Equal(t, "", item.Last)
}

func TestRowScansNestedFields(t *testing.T) {
	rows := queryRows(t, []string{"p_first", "p_last"}, []any{"Brett", "Jones"})

	var res struct {
		Item struct {
			First string `db:"p_first"`
			Last  string `db:"p_last"`
		}
	}

	require.NoError(t, sqlext.ScanRow(&res, rows))
	assert.Equal(t, "Brett", res.Item.First)
	assert.Equal(t, "Jones", res.Item.Last)
}

func TestRowScansNestedFieldsBySnakeCase(t *testing.T) {
	rows := queryRows(t, []string{"first_name", "last_name"}, []any{"Brett", "Jones"})

	var res struct {
		Item struct {
			FirstName string
			LastName  string
		}
	}

	require.NoError(t, sqlext.ScanRow(&res, rows))
	assert.Equal(t, "Brett", res.Item.FirstName)
	assert.Equal(t, "Jones", res.Item.LastName)
}

func TestRowsIgnoresDashTaggedFields(t *testing.T) {
	rows := queryRows(t, []string{"first", "last"},
		[]any{"Brett", "Jones"},
		[]any{"Fred", "Jones"},
	)

	var items []struct {
		First string
		Last  string `db:"-"`
	}

	require.NoError(t, sqlext.ScanRows(&items, rows))
	require.Len(t, items, 2)
	assert.Equal(t, "Brett", items[0].First)
	assert.Equal(t, "", items[0].Last)
	assert.Equal(t, "Fred", items[1].First)
	assert.Equal(t, "", items[1].Last)
}

func TestRowScansNestedPointerFields(t *testing.T) {
	rows := queryRows(t, []string{"id", "u_id", "u_name", "le_id", "le_name"},
		[]any{1, 7, "brett", 9, "fred"},
	)

	type user struct {
		ID   int
		Name string
	}
	var item struct {
		ID     int
		User   *user `db:"u_"`
		Editor *user `db:"le_"`
	}

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Equal(t, 1, item.ID)
	require.NotNil(t, item.User)
	assert.Equal(t, 7, item.User.ID)
	assert.Equal(t, "brett", item.User.Name)
	require.NotNil(t, item.Editor)
	assert.Equal(t, 9, item.Editor.ID)
	assert.Equal(t, "fred", item.Editor.Name)
}

func TestRowsScansNestedPointerFieldsAllocatingEachRow(t *testing.T) {
	rows := queryRows(t, []string{"id", "u_name"},
		[]any{1, "brett"},
		[]any{2, "fred"},
	)

	type user struct {
		Name string
	}
	var items []struct {
		ID   int
		User *user `db:"u_"`
	}

	require.NoError(t, sqlext.ScanRows(&items, rows))
	require.Len(t, items, 2)
	assert.Equal(t, "brett", items[0].User.Name)
	assert.Equal(t, "fred", items[1].User.Name)
}

func TestNestedStructPrefixAppliesToDBTags(t *testing.T) {
	rows := queryRows(t, []string{"c_id", "c_full_name"}, []any{3, "costco"})

	var item struct {
		Company struct {
			ID   int    `db:"id"`
			Name string `db:"full_name"`
		} `db:"c_"`
	}

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Equal(t, 3, item.Company.ID)
	assert.Equal(t, "costco", item.Company.Name)
}

func TestNestedStructPointerImplementingScannerIsLeaf(t *testing.T) {
	// a pointer to a struct that is a valid sql value (time.Time) stays a
	// single column: the scanner must not recurse into it, which would
	// allocate the nil pointer even when no column matches
	rows := queryRows(t, []string{"other"}, []any{"x"})

	var item struct {
		At    *time.Time
		Other string
	}

	require.NoError(t, sqlext.ScanRow(&item, rows))
	assert.Nil(t, item.At)
	assert.Equal(t, "x", item.Other)
}

func TestColumnsMapperSnakeCase(t *testing.T) {
	tests := map[string]string{
		"":              "",
		"id":            "id",
		"Name":          "name",
		"ID":            "id",
		"UUID":          "uuid",
		"UserID":        "user_id",
		"IsActive":      "is_active",
		"FirstName":     "first_name",
		"myCustomName":  "my_custom_name",
		"HTTPServer":    "http_server",
		"Address2Line":  "address2_line",
		"already_snake": "already_snake",
	}

	for in, out := range tests {
		assert.Equal(t, out, sqlext.ColumnsMapper(in), "ColumnsMapper(%q)", in)
	}
}
