package sqlext_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

// testDBCounter 保证同一个进程里每个测试用例拿到不重名的内存库。
var testDBCounter atomic.Int32

// newTestDB 打开一个独立的内存 SQLite 库，调用方关闭连接后内存库自动释放。
//
// 这里用带唯一名字的 shared cache 而不是 ":memory:"，是为了让连接池里的
// 多条连接看到同一个库：单条连接的 ":memory:" 库在换连接后就查不到建好的表。
func newTestDB() (*sql.DB, error) {
	name := fmt.Sprintf("sqlext_test_%d_%d", time.Now().UnixNano(), testDBCounter.Add(1))
	dsn := "file:" + name + "?mode=memory&cache=shared&_pragma=foreign_keys(ON)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open %q: %w", dsn, err)
	}

	// Ping 会真正建连，DSN 有问题时在这里就失败
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot connect to %q: %w", dsn, err)
	}
	return db, nil
}

// testDB 返回当前用例独占的内存库，用例结束后自动关闭。
// 建库失败直接失败用例而非跳过，避免驱动或 DSN 出问题时测试被静默略过。
func testDB(t testing.TB) *sql.DB {
	t.Helper()

	db, err := newTestDB()
	if err != nil {
		t.Fatalf("sqlext: %v", err)
	}

	// Cleanup 后注册的先执行，因此 queryRows 注册的 rows.Close 总在 db.Close 之前
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// queryRows 在内存 SQLite 上执行一条只由字面量组成的查询并返回真实的 *sql.Rows：
// cols 为结果集的列名，records 为各行取值（一行内的元素少于 cols 时以 NULL 补齐）。
// 例如 queryRows(t, []string{"id", "name"}, []any{1, "brett"}) 执行
//
//	SELECT 1 AS "id", 'brett' AS "name"
//
// 多行用 UNION ALL 串联，零行用 LIMIT 0 保留列信息。查询不依赖任何表，
// 因此不会与 example_scanner_test.go 重建的示例表相互影响。
func queryRows(t testing.TB, cols []string, records ...[]any) *sql.Rows {
	t.Helper()

	query := selectLiterals(cols, records)
	rows, err := testDB(t).Query(query)
	require.NoError(t, err, query)
	t.Cleanup(func() { rows.Close() })

	return rows
}

// selectLiterals 组装 queryRows 使用的 SELECT 语句。
func selectLiterals(cols []string, records [][]any) string {
	if len(records) == 0 {
		exprs := make([]string, len(cols))
		for i, col := range cols {
			exprs[i] = "NULL AS " + quoteIdent(col)
		}
		return "SELECT " + strings.Join(exprs, ", ") + " LIMIT 0"
	}

	selects := make([]string, 0, len(records))
	for _, record := range records {
		exprs := make([]string, len(cols))
		for i, col := range cols {
			var value any
			if i < len(record) {
				value = record[i]
			}

			exprs[i] = literal(value)
			if len(selects) == 0 {
				// UNION ALL 的结果列名取自第一个 SELECT
				exprs[i] += " AS " + quoteIdent(col)
			}
		}
		selects = append(selects, "SELECT "+strings.Join(exprs, ", "))
	}

	return strings.Join(selects, " UNION ALL ")
}

// quoteIdent 用双引号包裹标识符，双引号本身按 SQL 规则成对转义。
func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// literal 把测试值转换为 SQLite 字面量，覆盖 scanner 用例实际用到的类型。
// SQLite 用两个单引号转义单引号，反斜杠不是转义字符。
func literal(value any) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	}

	panic(fmt.Sprintf("sqlext test: unsupported literal type %T", value))
}
