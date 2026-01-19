package sqlutil

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// Test struct with various MySQL types
type TestAllTypes struct {
	ID          int64           `db:"id"`
	TinyInt     int8            `db:"tiny_int"`
	SmallInt    int16           `db:"small_int"`
	MediumInt   int32           `db:"medium_int"`
	BigInt      int64           `db:"big_int"`
	Float       float32         `db:"float_val"`
	Double      float64         `db:"double_val"`
	Decimal     string          `db:"decimal_val"`
	Varchar     string          `db:"varchar_val"`
	Text        string          `db:"text_val"`
	Blob        []byte          `db:"blob_val"`
	DateTime    time.Time       `db:"datetime_val"`
	Bool        bool            `db:"bool_val"`
	JSON        string          `db:"json_val"`
	NullString  sql.NullString  `db:"null_string"`
	NullInt64   sql.NullInt64   `db:"null_int64"`
	NullFloat64 sql.NullFloat64 `db:"null_float64"`
	NullBool    sql.NullBool    `db:"null_bool"`
	NullTime    sql.NullTime    `db:"null_time"`
}

func TestQuery_AllTypes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Now().Truncate(time.Second)
	columns := []string{
		"id", "tiny_int", "small_int", "medium_int", "big_int",
		"float_val", "double_val", "decimal_val", "varchar_val", "text_val",
		"blob_val", "datetime_val", "bool_val", "json_val",
		"null_string", "null_int64", "null_float64", "null_bool", "null_time",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		int64(1), int8(127), int16(32767), int32(8388607), int64(9223372036854775807),
		float32(3.14), float64(3.14159265359), "12345.67", "hello", "long text",
		[]byte("binary data"), now, true, `{"key":"value"}`,
		"valid string", int64(42), float64(99.99), true, now,
	)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	ctx := context.Background()
	results, err := Query[TestAllTypes](ctx, db, "SELECT * FROM test_table")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.ID != 1 {
		t.Errorf("expected ID=1, got %d", r.ID)
	}
	if r.TinyInt != 127 {
		t.Errorf("expected TinyInt=127, got %d", r.TinyInt)
	}
	if r.SmallInt != 32767 {
		t.Errorf("expected SmallInt=32767, got %d", r.SmallInt)
	}
	if r.Varchar != "hello" {
		t.Errorf("expected Varchar='hello', got '%s'", r.Varchar)
	}
	if r.Bool != true {
		t.Errorf("expected Bool=true, got %v", r.Bool)
	}
	if !r.NullString.Valid || r.NullString.String != "valid string" {
		t.Errorf("expected NullString='valid string', got %v", r.NullString)
	}
	if !r.NullInt64.Valid || r.NullInt64.Int64 != 42 {
		t.Errorf("expected NullInt64=42, got %v", r.NullInt64)
	}
}

func TestQuery_NullValues(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	columns := []string{"null_string", "null_int64", "null_float64", "null_bool", "null_time"}
	rows := sqlmock.NewRows(columns).AddRow(nil, nil, nil, nil, nil)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	type NullOnly struct {
		NullString  sql.NullString  `db:"null_string"`
		NullInt64   sql.NullInt64   `db:"null_int64"`
		NullFloat64 sql.NullFloat64 `db:"null_float64"`
		NullBool    sql.NullBool    `db:"null_bool"`
		NullTime    sql.NullTime    `db:"null_time"`
	}

	ctx := context.Background()
	results, err := Query[NullOnly](ctx, db, "SELECT * FROM null_table")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.NullString.Valid {
		t.Errorf("expected NullString to be NULL, got %v", r.NullString)
	}
	if r.NullInt64.Valid {
		t.Errorf("expected NullInt64 to be NULL, got %v", r.NullInt64)
	}
	if r.NullFloat64.Valid {
		t.Errorf("expected NullFloat64 to be NULL, got %v", r.NullFloat64)
	}
	if r.NullBool.Valid {
		t.Errorf("expected NullBool to be NULL, got %v", r.NullBool)
	}
	if r.NullTime.Valid {
		t.Errorf("expected NullTime to be NULL, got %v", r.NullTime)
	}
}

func TestQueryOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	type Simple struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		ctx := context.Background()
		result, err := QueryOne[Simple](ctx, db, "SELECT * FROM simple WHERE id = ?", 1)
		if err != nil {
			t.Fatalf("QueryOne failed: %v", err)
		}
		if result.ID != 1 || result.Name != "test" {
			t.Errorf("unexpected result: %+v", result)
		}
	})

	t.Run("not found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"})
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		ctx := context.Background()
		_, err := QueryOne[Simple](ctx, db, "SELECT * FROM simple WHERE id = ?", 999)
		if err != sql.ErrNoRows {
			t.Errorf("expected sql.ErrNoRows, got %v", err)
		}
	})
}

func TestExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 5))

	ctx := context.Background()
	affected, err := Exec(ctx, db, "UPDATE users SET status = ? WHERE active = ?", 0, true)
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if affected != 5 {
		t.Errorf("expected 5 rows affected, got %d", affected)
	}
}

func TestExecResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(123, 1))

	ctx := context.Background()
	result, err := ExecResult(ctx, db, "INSERT INTO users (name) VALUES (?)", "test")
	if err != nil {
		t.Fatalf("ExecResult failed: %v", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("LastInsertId failed: %v", err)
	}
	if lastID != 123 {
		t.Errorf("expected LastInsertId=123, got %d", lastID)
	}
}

func TestQuery_MultipleRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	type Item struct {
		ID    int    `db:"id"`
		Value string `db:"value"`
	}

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow(1, "first").
		AddRow(2, "second").
		AddRow(3, "third")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	ctx := context.Background()
	results, err := Query[Item](ctx, db, "SELECT * FROM items")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	expected := []struct {
		id    int
		value string
	}{
		{1, "first"},
		{2, "second"},
		{3, "third"},
	}

	for i, exp := range expected {
		if results[i].ID != exp.id || results[i].Value != exp.value {
			t.Errorf("row %d: expected {%d, %s}, got {%d, %s}",
				i, exp.id, exp.value, results[i].ID, results[i].Value)
		}
	}
}

func TestQuery_UnmappedColumn(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	type Partial struct {
		ID int `db:"id"`
		// extra_col is not mapped
	}

	rows := sqlmock.NewRows([]string{"id", "extra_col"}).AddRow(1, "ignored")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	ctx := context.Background()
	results, err := Query[Partial](ctx, db, "SELECT * FROM partial_table")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 || results[0].ID != 1 {
		t.Errorf("unexpected result: %+v", results)
	}
}

func TestQuery_IgnoreTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	type WithIgnore struct {
		ID      int    `db:"id"`
		Ignored string `db:"-"`
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	ctx := context.Background()
	results, err := Query[WithIgnore](ctx, db, "SELECT id FROM ignore_table")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 || results[0].ID != 42 {
		t.Errorf("unexpected result: %+v", results)
	}
	if results[0].Ignored != "" {
		t.Errorf("expected Ignored to be empty, got '%s'", results[0].Ignored)
	}
}

func TestQuery_NullToNonPointerTypes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Struct with non-pointer types (string, int, float, bool)
	type NonPointer struct {
		ID      int     `db:"id"`
		Name    string  `db:"name"`
		Age     int     `db:"age"`
		Score   float64 `db:"score"`
		Active  bool    `db:"active"`
		TinyVal int8    `db:"tiny_val"`
	}

	// Return NULL for all columns except ID
	rows := sqlmock.NewRows([]string{"id", "name", "age", "score", "active", "tiny_val"}).
		AddRow(1, nil, nil, nil, nil, nil)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	ctx := context.Background()
	results, err := Query[NonPointer](ctx, db, "SELECT * FROM non_pointer_table")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	// ID should be 1
	if r.ID != 1 {
		t.Errorf("expected ID=1, got %d", r.ID)
	}
	// NULL values should be zero values
	if r.Name != "" {
		t.Errorf("expected Name='', got '%s'", r.Name)
	}
	if r.Age != 0 {
		t.Errorf("expected Age=0, got %d", r.Age)
	}
	if r.Score != 0 {
		t.Errorf("expected Score=0, got %f", r.Score)
	}
	if r.Active != false {
		t.Errorf("expected Active=false, got %v", r.Active)
	}
	if r.TinyVal != 0 {
		t.Errorf("expected TinyVal=0, got %d", r.TinyVal)
	}
}

func TestQuery_MixedNullAndValues(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	type Mixed struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	// Row 1: all values, Row 2: name is NULL, Row 3: email is NULL
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "Alice", "alice@example.com").
		AddRow(2, nil, "bob@example.com").
		AddRow(3, "Charlie", nil)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	ctx := context.Background()
	results, err := Query[Mixed](ctx, db, "SELECT * FROM mixed_table")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// Row 1
	if results[0].ID != 1 || results[0].Name != "Alice" || results[0].Email != "alice@example.com" {
		t.Errorf("row 0 unexpected: %+v", results[0])
	}
	// Row 2 - Name is NULL, should be ""
	if results[1].ID != 2 || results[1].Name != "" || results[1].Email != "bob@example.com" {
		t.Errorf("row 1 unexpected: %+v", results[1])
	}
	// Row 3 - Email is NULL, should be ""
	if results[2].ID != 3 || results[2].Name != "Charlie" || results[2].Email != "" {
		t.Errorf("row 2 unexpected: %+v", results[2])
	}
}
