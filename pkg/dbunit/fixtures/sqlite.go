package fixtures

import (
	"database/sql"
	"fmt"
)

type loadFunction func(tx *sql.Tx) error

// sqliteHelper SQLite 数据库辅助工具，实现表名查询、关键字引用、外键控制等
type sqliteHelper struct {
	tables []string
}

func (h *sqliteHelper) init(db *sql.DB) error {
	var err error
	h.tables, err = h.tableNames(db)
	if err != nil {
		return err
	}

	return nil
}

func (*sqliteHelper) quoteKeyword(str string) string {
	return fmt.Sprintf("`%s`", str)
}

// databaseName 返回数据库名称（文件路径），包含 test_ 前缀以满足 EnsureTestDatabase 检查
func (*sqliteHelper) databaseName(q *sql.DB) (string, error) {
	// PRAGMA database_list 返回 3 列: seq, name, file
	var seq int
	var name, file string
	err := q.QueryRow("PRAGMA database_list").Scan(&seq, &name, &file)
	if err != nil {
		return "", err
	}
	return file, nil
}

// tableNames 查询 SQLite 数据库中的所有用户表
func (*sqliteHelper) tableNames(q *sql.DB) ([]string, error) {
	query := `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table'
		  AND name NOT LIKE 'sqlite_%'
		ORDER BY name;
	`

	rows, err := q.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

// disableReferentialIntegrity 临时关闭外键检查，执行 loadFn 后重新开启
func (*sqliteHelper) disableReferentialIntegrity(db *sql.DB, loadFn loadFunction) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err = tx.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return err
	}

	err = loadFn(tx)
	if _, err2 := tx.Exec("PRAGMA foreign_keys = ON"); err2 != nil {
		if err == nil {
			err = err2
		}
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}
