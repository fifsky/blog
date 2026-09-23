package sqlext

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"
)

// Builder builds a SQL query together with the positional arguments (?) used by
// database/sql. The zero value is ready to use and every method returns the
// receiver so the calls can be chained.
//
//	q := sqlext.NewBuilder().
//		Select("api.id, api.api_name").
//		From("api LEFT JOIN system_menu_api ON system_menu_api.api_id = api.id").
//		Where("api.service_id = ? AND api.status = ?", serviceID, status).
//		WhereIf(keyword != "", "(api.api_name LIKE ? or api.api_url LIKE ?)", kw, kw).
//		GroupBy("api.id").
//		Having("COUNT(*) > ?", 1).
//		OrderBy("api.created_at desc").
//		Limit(10).
//		Offset(20)
//
//	rows, err := db.Query(q.SQL(), q.Args()...)
//	// SELECT api.id, api.api_name
//	//   FROM api LEFT JOIN system_menu_api ON system_menu_api.api_id = api.id
//	//   WHERE api.service_id = ? AND api.status = ? AND (api.api_name LIKE ? or api.api_url LIKE ?)
//	//   GROUP BY api.id HAVING COUNT(*) > ? ORDER BY api.created_at desc LIMIT 10 OFFSET 20
//
// Joins are written directly inside From and Args always returns the arguments
// in the order of the placeholders inside SQL, no matter in which order the
// methods have been called.
type Builder struct {
	columns    []string
	selectArgs []any
	from       []string
	fromArgs   []any
	wheres     []condition
	groups     []string
	havings    []condition
	orders     []string
	orderArgs  []any
	limit      int
	offset     int
}

// condition is a single expression together with its arguments and the
// connector used to join it with the previous expression.
type condition struct {
	connector string
	expr      string
	args      []any
}

// NewBuilder returns an empty Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// Select sets the selected columns e.g. Select("api.id, api.api_name"). Empty
// values are ignored and * is selected when nothing is set. The expression
// may contain ? placeholders (e.g. expressions with LOCATE(?)) which have to
// be matched by args. Multiple calls are joined with commas.
func (b *Builder) Select(columns string, args ...any) *Builder {
	columns = strings.TrimSpace(columns)
	if columns == "" {
		return b
	}

	b.columns = append(b.columns, columns)
	b.selectArgs = append(b.selectArgs, args...)

	return b
}

// From sets the FROM clause, joins are written directly e.g.
// From("api LEFT JOIN company ON company.id = api.company_id"). The expression
// may contain ? placeholders which have to be matched by args. Multiple calls
// are joined with commas.
func (b *Builder) From(expression string, args ...any) *Builder {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return b
	}

	b.from = append(b.from, expression)
	b.fromArgs = append(b.fromArgs, args...)

	return b
}

// Where appends a condition joined with the previous one using AND. The
// condition may contain any number of ? placeholders which have to be matched
// by args. Empty conditions are ignored so dynamic conditions are safe, and a
// slice argument is expanded into one placeholder per element, so an IN list
// is written with a single placeholder:
//
//	q.Where("doc_id = ?", docID)
//	q.Where("id IN (?)", ids)
func (b *Builder) Where(condition string, args ...any) *Builder {
	b.wheres = appendCondition(b.wheres, "AND", condition, args)

	return b
}

// WhereIf is identical to Where, but the condition and its arguments are only
// appended when cond is true.
func (b *Builder) WhereIf(cond bool, condition string, args ...any) *Builder {
	if !cond {
		return b
	}

	b.wheres = appendCondition(b.wheres, "AND", condition, args)

	return b
}

// OrWhere is identical to Where, but the condition is joined with the previous
// one using OR. When it is the first condition the connector is ignored.
func (b *Builder) OrWhere(condition string, args ...any) *Builder {
	b.wheres = appendCondition(b.wheres, "OR", condition, args)

	return b
}

// OrWhereIf is identical to OrWhere, but the condition and its arguments are
// only appended when cond is true.
func (b *Builder) OrWhereIf(cond bool, condition string, args ...any) *Builder {
	if !cond {
		return b
	}

	b.wheres = appendCondition(b.wheres, "OR", condition, args)

	return b
}

// GroupBy sets the GROUP BY columns e.g. GroupBy("api.id, api.status"). Empty
// values are ignored, multiple calls are joined with commas.
func (b *Builder) GroupBy(columns string) *Builder {
	b.groups = appendStrings(b.groups, columns)

	return b
}

// Group is an alias for GroupBy.
func (b *Builder) Group(columns string) *Builder {
	return b.GroupBy(columns)
}

// Having appends a condition to the HAVING clause. Conditions are joined with
// AND and empty conditions are ignored.
func (b *Builder) Having(condition string, args ...any) *Builder {
	b.havings = appendCondition(b.havings, "AND", condition, args)

	return b
}

// HavingIf is identical to Having, but the condition and its arguments are
// only appended when cond is true.
func (b *Builder) HavingIf(cond bool, condition string, args ...any) *Builder {
	if !cond {
		return b
	}

	b.havings = appendCondition(b.havings, "AND", condition, args)

	return b
}

// OrHaving is identical to Having, but the condition is joined with the
// previous one using OR. When it is the first condition the connector is
// ignored.
func (b *Builder) OrHaving(condition string, args ...any) *Builder {
	b.havings = appendCondition(b.havings, "OR", condition, args)

	return b
}

// OrHavingIf is identical to OrHaving, but the condition and its arguments are
// only appended when cond is true.
func (b *Builder) OrHavingIf(cond bool, condition string, args ...any) *Builder {
	if !cond {
		return b
	}

	b.havings = appendCondition(b.havings, "OR", condition, args)

	return b
}

// OrderBy sets the ORDER BY expressions e.g. OrderBy("created_at desc, id").
// Empty values are ignored, multiple calls are joined with commas. The
// expressions may contain ? placeholders (e.g. fulltext relevance ordering
// MATCH(...) AGAINST(?)) which have to be matched by args.
func (b *Builder) OrderBy(orders string, args ...any) *Builder {
	b.orders = appendStrings(b.orders, orders)
	b.orderArgs = append(b.orderArgs, args...)

	return b
}

// Order is an alias for OrderBy.
func (b *Builder) Order(orders string, args ...any) *Builder {
	return b.OrderBy(orders, args...)
}

// Limit limits the number of returned rows. A value smaller than or equal to
// zero removes the limit.
func (b *Builder) Limit(n int) *Builder {
	b.limit = n

	return b
}

// Offset skips the first n rows. A value smaller than or equal to zero removes
// the offset. Note that MySQL requires a limit to be set as well.
func (b *Builder) Offset(n int) *Builder {
	b.offset = n

	return b
}

// SQL returns the built SQL query.
func (b *Builder) SQL() string {
	var sb strings.Builder

	if len(b.columns) > 0 {
		sb.WriteString("SELECT ")
		sb.WriteString(strings.Join(b.columns, ", "))
	} else {
		sb.WriteString("SELECT *")
	}

	if len(b.from) > 0 {
		sb.WriteString(" FROM ")
		sb.WriteString(strings.Join(b.from, ", "))
	}

	if where := renderConditions(b.wheres); where != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(where)
	}

	if len(b.groups) > 0 {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(strings.Join(b.groups, ", "))
	}

	if having := renderConditions(b.havings); having != "" {
		sb.WriteString(" HAVING ")
		sb.WriteString(having)
	}

	if len(b.orders) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(b.orders, ", "))
	}

	if b.limit > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(b.limit))
	}

	if b.offset > 0 {
		sb.WriteString(" OFFSET ")
		sb.WriteString(strconv.Itoa(b.offset))
	}

	return sb.String()
}

// Args returns the positional arguments matching the ? placeholders of SQL.
// Arguments are returned in the order in which the placeholders appear in SQL.
func (b *Builder) Args() []any {
	var args []any

	args = append(args, b.selectArgs...)
	args = append(args, b.fromArgs...)
	for _, c := range b.wheres {
		args = append(args, c.args...)
	}
	for _, c := range b.havings {
		args = append(args, c.args...)
	}
	args = append(args, b.orderArgs...)

	return args
}

// Build returns the query built by SQL together with its arguments.
func (b *Builder) Build() (string, []any) {
	return b.SQL(), b.Args()
}

// String implements fmt.Stringer and is an alias for SQL so a builder can be
// used directly in a formatted query.
func (b *Builder) String() string {
	return b.SQL()
}

// appendCondition appends expr to conds, empty expressions are ignored.
// Slice arguments are expanded into one placeholder per element, see expandArgs.
func appendCondition(conds []condition, connector, expr string, args []any) []condition {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return conds
	}

	expr, args = expandArgs(expr, args)

	return append(conds, condition{connector: connector, expr: expr, args: args})
}

// expandArgs rewrites expr so that every ? placeholder whose argument is a
// slice gets one placeholder per element, returning the expression together
// with the flattened arguments:
//
//	expandArgs("id IN (?)", []any{[]int{1, 2, 3}})
//	// "id IN (?,?,?)", []any{1, 2, 3}
//
// An empty slice becomes NULL: an empty list is not valid SQL and IN (NULL)
// matches nothing, which is the expected result of a list filter without
// values. Slices that database/sql treats as a single value are not expanded,
// see sliceValues.
func expandArgs(expr string, args []any) (string, []any) {
	if !hasSliceArg(args) {
		return expr, args
	}

	var sb strings.Builder
	sb.Grow(len(expr) + 2*len(args))

	expanded := make([]any, 0, len(args))
	next := 0

	for i := 0; i < len(expr); i++ {
		if expr[i] != '?' || next == len(args) {
			sb.WriteByte(expr[i])
			continue
		}

		value := args[next]
		next++

		values, ok := sliceValues(value)
		if !ok {
			sb.WriteByte('?')
			expanded = append(expanded, value)
			continue
		}

		if len(values) == 0 {
			sb.WriteString("NULL")
			continue
		}

		for j, v := range values {
			if j > 0 {
				sb.WriteByte(',')
			}
			sb.WriteByte('?')
			expanded = append(expanded, v)
		}
	}

	return sb.String(), expanded
}

// hasSliceArg reports whether args contains a slice that has to be expanded.
func hasSliceArg(args []any) bool {
	for _, arg := range args {
		if _, ok := sliceValues(arg); ok {
			return true
		}
	}

	return false
}

// sliceValues returns the elements of a slice or array argument. It reports
// false for a nil or non slice value, and for slices that have to stay a
// single argument: []byte is a binary value for database/sql, and a type
// implementing driver.Valuer (e.g. a comma joined id list) converts itself
// into one value.
func sliceValues(arg any) ([]any, bool) {
	v := reflect.ValueOf(arg)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return nil, false
	}

	if v.Type().Elem().Kind() == reflect.Uint8 || v.Type().Implements(reflect.TypeFor[driver.Valuer]()) {
		return nil, false
	}

	values := make([]any, v.Len())
	for i := range values {
		values[i] = v.Index(i).Interface()
	}

	return values, true
}

// appendStrings appends the non empty values to values.
func appendStrings(values []string, items ...string) []string {
	for _, item := range items {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}

	return values
}

// renderConditions joins conditions with their connectors.
func renderConditions(conds []condition) string {
	var sb strings.Builder

	for i, c := range conds {
		if i > 0 {
			sb.WriteString(" ")
			sb.WriteString(c.connector)
			sb.WriteString(" ")
		}
		sb.WriteString(c.expr)
	}

	return sb.String()
}
