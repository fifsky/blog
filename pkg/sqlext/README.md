## Examples

### Multiple Rows

```go
db, err := sql.Open("sqlite3", "database.sqlite")
rows, err := db.Query("SELECT * FROM persons")
defer rows.Close()

var persons []Person
err = sqlext.ScanRows(&persons, rows)

fmt.Printf("%#v", persons)
// []Person{
//    {ID: 1, Name: "brett"},
//    {ID: 2, Name: "fred"},
//    {ID: 3, Name: "stacy"},
// }
```

### Multiple rows of primitive type

```go
rows, err := db.Query("SELECT name FROM persons")
defer rows.Close()

var names []string
err = sqlext.ScanRows(&names, rows)

fmt.Printf("%#v", names)
// []string{
//    "brett",
//    "fred",
//    "stacy",
// }
```

### Single row

```go
rows, err := db.Query("SELECT * FROM persons where name = 'brett' LIMIT 1")
defer rows.Close()

var person Person
err = sqlext.ScanRow(&person, rows)

fmt.Printf("%#v", person)
// Person{ ID: 1, Name: "brett" }
```

### Scalar value

```go
rows, err := db.Query("SELECT age FROM persons where name = 'brett' LIMIT 1")
defer rows.Close()

var age int8
err = sqlext.ScanRow(&age, rows)

fmt.Printf("%d", age)
// 100
```

### Nested Struct Fields

The fields of nested structs are matched with the columns exactly like top level fields:

```go
rows, err := db.Query(`
	SELECT person.id AS person_id, person.name AS person_name, company.name AS company_name
	FROM person
	JOIN company on company.id = person.company_id
	LIMIT 1
`)
defer rows.Close()

var person struct {
	ID      int    `db:"person_id"`
	Name    string `db:"person_name"`
	Company struct {
		Name string `db:"company_name"`
	}
}

err = sqlext.ScanRow(&person, rows)

err = json.NewEncoder(os.Stdout).Encode(&person)
// Output:
// {"ID":1,"Name":"brett","Company":{"Name":"costco"}}
```

Nested fields may also be pointers to structs: a nil pointer is allocated in
place so the row scans into it. A `db` tag on a nested struct (or pointer to
struct) is used as a **column prefix** for all of its fields. This is how self
joins are scanned — alias the repeated columns in the query and prefix the
nested models:

```go
rows, err := db.Query(`
	SELECT a.id, a.title,
	       u.id AS u_id, u.name AS u_name,
	       le.id AS le_id, le.name AS le_name
	FROM articles a
	LEFT JOIN users u ON a.user_id = u.id
	LEFT JOIN users le ON a.last_user_id = le.id
`)
defer rows.Close()

var article struct {
	ID         int
	Title      string
	User       *User `db:"u_"`   // matches u_id, u_name...
	LastEditor *User `db:"le_"`  // matches le_id, le_name...
}

err = sqlext.ScanRows(&article, rows)
```

Structs that are valid sql values themselves (`time.Time`, `driver.Valuer`
implementations) are scanned as a single column and never recursed into.

### QueryRow and Query

`QueryRow` and `Query` run the query and close the rows themselves, so callers
only deal with the scanned value. `QueryRow` scans the first row into `T` and
returns `sql.ErrNoRows` when there is no row; `Query` scans every row into
`[]T` and returns an empty (non nil) slice for an empty result. Both take a
`Queryer`, which `*sql.DB` and `*sql.Tx` satisfy.

```go
q := sqlext.NewBuilder().
        Select("id, name").
        From("persons").
        Where("age > ?", 100)

persons, err := sqlext.Query[Person](ctx, db, q.SQL(), q.Args()...)
person, err := sqlext.QueryRow[Person](ctx, db, q.SQL(), q.Args()...)
```

### Column Names

The column name of a field is determined by:

- the `db` struct tag when it is set, e.g. `Name string \`db:"nama"\`` is matched with the `nama` column
- the snake_case form of the field name when the tag is missing, e.g. `FirstName` is matched with the `first_name` column
- fields tagged with `db:"-"` are ignored and never matched

You can override the snake_case behavior by setting `ColumnsMapper` to a custom function.

The tag has to match the column name returned by the driver, not the expression
written in the query: MySQL strips the table qualifier and lets duplicated names
collide, so `SELECT person.id, person.name, company.name` comes back as the
columns `id`, `name` and `name`. Use an alias (`AS person_id`) whenever a
qualified or repeated column is selected, otherwise it cannot be matched.

### Builder

`Builder` builds a complete SQL query with the positional arguments used by `database/sql`.

```go
q := sqlext.NewBuilder().
        Select("api.id, api.api_name").
        From("api LEFT JOIN system_menu_api ON system_menu_api.api_id = api.id").
        Where("api.service_id = ? AND api.status = ? and api.api_type = ?", request.ServiceID, constant.StatusOnline, request.ApiType).
        WhereIf(request.Keyword != "", "(api.api_name LIKE ? or api.api_url LIKE ?)", kw, kw).
        WhereIf(request.NotInMenuId > 0, "api.id NOT IN (SELECT api_id FROM system_menu_api WHERE menu_id = ?)", request.NotInMenuId).
        GroupBy("api.id").
        Having("COUNT(*) > ?", 1).
        OrderBy("api.created_at desc").
        Limit(20).
        Offset(0)

rows, err := db.Query(q.SQL(), q.Args()...)
// SELECT api.id, api.api_name FROM api LEFT JOIN system_menu_api ON system_menu_api.api_id = api.id
//   WHERE api.service_id = ? AND api.status = ? and api.api_type = ? AND (api.api_name LIKE ? or api.api_url LIKE ?) AND api.id NOT IN (SELECT api_id FROM system_menu_api WHERE menu_id = ?)
//   GROUP BY api.id HAVING COUNT(*) > ? ORDER BY api.created_at desc LIMIT 20
```

| method                                      | description                                                |
| ------------------------------------------- | ---------------------------------------------------------- |
| `Select(columns string, args ...)`          | selected columns, defaults to `*`; `?` placeholders are supported (e.g. `LOCATE(?, col)`) |
| `From(expression string, args ...)`         | FROM clause, joins are written inline and `?` is supported |
| `Where(expr string, args ...)`              | condition joined with `AND`, empty expressions are ignored |
| `WhereIf(cond bool, expr string, args ...)` | condition appended only when `cond` is true                |
| `OrWhere/OrWhereIf`                         | same as above but joined with `OR`                         |
| `GroupBy(columns string)` / `Group`         | `GROUP BY`                                                 |
| `Having/HavingIf/OrHaving/OrHavingIf`       | `HAVING` conditions                                        |
| `OrderBy(orders string, args ...)` / `Order`| `ORDER BY`; `?` placeholders are supported (e.g. fulltext `AGAINST(?)` ordering) |
| `Limit(n int)` / `Offset(n int)`            | non positive values are ignored                            |
| `SQL()` / `Args()` / `Build()`              | the query, the arguments and both together                 |

The zero value is ready to use, every method returns the receiver so calls can be chained, and `Args` always follows the order of the placeholders inside `SQL`, no matter in which order the methods were called.

A slice argument is expanded into one placeholder per element, so an `IN` list
is written with a single placeholder:

```go
q := sqlext.NewBuilder().
        Select("id, name").
        From("persons").
        Where("status = ?", 1).
        WhereIf(len(ids) > 0, "id IN (?)", ids)
// SELECT id, name FROM persons WHERE status = ? AND id IN (?,?,?)
```

An empty slice becomes `NULL`, because an empty list is not valid SQL and
`IN (NULL)` matches nothing. Slices that `database/sql` uses as a single value
are never expanded: `[]byte` (binary data) and types implementing
`driver.Valuer`.

For statements the Builder cannot generate, e.g. a `DELETE` with an `IN` list,
`In` returns the placeholders and the arguments of the clause:

```go
placeholders, args := sqlext.In(ids)
// placeholders = "?,?,?"   args = []any{1, 2, 3}
db.Exec("delete from persons where id in ("+placeholders+")", args...)
```

An empty slice returns no placeholders and no arguments: the caller is
responsible for skipping the statement, as an empty `IN` list is not valid SQL.

## Closing Rows

`ScanRow` and `ScanRows` never close the rows they are given. The caller is responsible for closing them:

```go
rows, err := db.Query("SELECT * FROM persons")
if err != nil {
        return err
}
defer rows.Close()

var persons []Person
if err := sqlext.ScanRows(&persons, rows); err != nil {
        return err
}
```
