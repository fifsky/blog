package sqlext_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"app/pkg/sqlext"
)

// exampleTable is dropped and recreated by every example
const exampleTable = "sqlext_example_person"

// exampleDB 打开一个独立的内存 SQLite 库并重建 exampleTable：
//
//	id | name  | company_name
//	1  | brett | costco
//	2  | fred  | costco
//	3  | NULL  | NULL
//
// 返回的连接由调用方关闭，关闭后内存库随之释放。
func exampleDB() *sql.DB {
	db, err := newTestDB()
	if err != nil {
		panic(fmt.Errorf("sqlext example: %w", err))
	}

	queries := []string{
		`DROP TABLE IF EXISTS ` + exampleTable,
		`CREATE TABLE ` + exampleTable + ` (
			id INTEGER NOT NULL PRIMARY KEY,
			name TEXT,
			company_name TEXT
		)`,
		`INSERT INTO ` + exampleTable + ` (id, name, company_name) VALUES (1, 'brett', 'costco')`,
		`INSERT INTO ` + exampleTable + ` (id, name, company_name) VALUES (2, 'fred', 'costco')`,
		`INSERT INTO ` + exampleTable + ` (id, name, company_name) VALUES (3, NULL, NULL)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			db.Close()
			panic(fmt.Errorf("sqlext example: %q: %w", query, err))
		}
	}

	return db
}

func ExampleScanRow() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM " + exampleTable + " ORDER BY id ASC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var person struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	err = sqlext.ScanRow(&person, rows)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(&person)
	if err != nil {
		panic(err)
	}
	// Output:
	// {"ID":1,"Name":"brett"}
}

func ExampleScanRow_nested() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, company_name FROM " + exampleTable + " ORDER BY id ASC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// nested struct fields are matched with columns just like top level fields
	var person struct {
		Info struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
		}
		Company struct {
			Name string `db:"company_name"`
		}
	}

	err = sqlext.ScanRow(&person, rows)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(&person)
	if err != nil {
		panic(err)
	}
	// Output:
	// {"Info":{"ID":1,"Name":"brett"},"Company":{"Name":"costco"}}
}

func ExampleScanRow_ignoredField() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM " + exampleTable + " ORDER BY id ASC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var person struct {
		ID   int `db:"-"`
		Name string
	}

	err = sqlext.ScanRow(&person, rows)
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(os.Stdout).Encode(&person)
	// Output:
	// {"ID":0,"Name":"brett"}
}

func ExampleScanRow_pointer() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM " + exampleTable + " WHERE id = 3 LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var person struct {
		ID   int
		Name *string `db:"name"`
	}

	err = sqlext.ScanRow(&person, rows)
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(os.Stdout).Encode(&person)
	// Output:
	// {"ID":3,"Name":null}
}

func ExampleScanRow_pointerType() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM " + exampleTable + " WHERE id = 3 LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	type NullableString *string
	var person struct {
		ID   int
		Name NullableString `db:"name"`
	}

	err = sqlext.ScanRow(&person, rows)
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(os.Stdout).Encode(&person)
	// Output:
	// {"ID":3,"Name":null}
}

func ExampleScanRow_scalar() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT name FROM " + exampleTable + " ORDER BY id ASC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var name string

	err = sqlext.ScanRow(&name, rows)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q", name)
	// Output:
	// "brett"
}

func ExampleScanRows() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM " + exampleTable + " ORDER BY id ASC")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var persons []struct {
		ID   int     `db:"id"`
		Name *string `db:"name"`
	}

	err = sqlext.ScanRows(&persons, rows)
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(os.Stdout).Encode(&persons)
	// Output:
	// [{"ID":1,"Name":"brett"},{"ID":2,"Name":"fred"},{"ID":3,"Name":null}]
}

func ExampleScanRows_snakeCase() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM " + exampleTable + " ORDER BY id ASC")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// fields without a db tag are matched by the snake_case form of their name
	var persons []struct {
		ID   int
		Name *string
	}

	err = sqlext.ScanRows(&persons, rows)
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(os.Stdout).Encode(&persons)
	// Output:
	// [{"ID":1,"Name":"brett"},{"ID":2,"Name":"fred"},{"ID":3,"Name":null}]
}

func ExampleScanRows_primitive() {
	db := exampleDB()
	defer db.Close()

	rows, err := db.Query("SELECT name FROM " + exampleTable + " WHERE name IS NOT NULL ORDER BY id ASC")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var names []string
	err = sqlext.ScanRows(&names, rows)
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(os.Stdout).Encode(&names)
	// Output:
	// ["brett","fred"]
}
