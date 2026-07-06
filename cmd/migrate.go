package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"app/config"
	"app/testutil"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"

	"github.com/urfave/cli/v3"
)

// migrateTables 定义迁移表及其列顺序（与 schema.sql 对应）
var migrateTables = []struct {
	name    string
	columns []string
}{
	{"cates", []string{"id", "name", "desc", "domain", "created_at", "updated_at"}},
	{"links", []string{"id", "name", "url", "desc", "status", "created_at", "updated_at"}},
	{"moods", []string{"id", "content", "user_id", "created_at", "updated_at"}},
	{"options", []string{"id", "option_key", "option_value"}},
	{"users", []string{"id", "name", "password", "nick_name", "email", "status", "type", "totp_secret", "openid", "created_at", "updated_at"}},
	{"posts", []string{"id", "cate_id", "type", "user_id", "title", "url", "content", "view_num", "tags", "status", "created_at", "updated_at"}},
	{"comments", []string{"id", "post_id", "pid", "name", "reply_name", "email", "website", "content", "ip", "created_at"}},
	{"reminds", []string{"id", "cron", "content", "status", "next_time", "created_at", "updated_at"}},
	{"regions", []string{"region_id", "parent_id", "level", "region_name", "longitude", "latitude", "pinyin", "az_no"}},
	{"guestbook", []string{"id", "name", "content", "ip", "top", "created_at"}},
	{"footprints", []string{"id", "name", "description", "longitude", "latitude", "date", "marker_color", "categories", "url", "url_label", "photos", "created_at", "updated_at"}},
}

// NewMigrate 创建 MySQL→SQLite 数据迁移命令
func NewMigrate() *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "migrate data from MySQL to SQLite",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "mysql-dsn",
				Usage: "MySQL source DSN, eg: root:123456@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
			},
			&cli.StringFlag{
				Name:  "sqlite-path",
				Usage: "SQLite destination file path, eg: storage/blog.db",
				Value: "storage/blog.db",
			},
		},
		Action: func(ctx context.Context, cli *cli.Command) error {
			mysqlDSN := cli.String("mysql-dsn")
			if mysqlDSN == "" {
				// 尝试从 config.yml 读取（兼容旧 MySQL 配置）
				conf := config.New()
				mysqlDSN = fmt.Sprintf("root:123456@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai")
				_ = conf
			}

			sqlitePath := cli.String("sqlite-path")
			if sqlitePath == "" {
				sqlitePath = "storage/blog.db"
			}

			log.Printf("[migrate] MySQL DSN: %s", mysqlDSN)
			log.Printf("[migrate] SQLite path: %s", sqlitePath)

			// 连接 MySQL
			mysqlDB, err := sql.Open("mysql", mysqlDSN)
			if err != nil {
				return fmt.Errorf("connect MySQL error: %w", err)
			}
			defer mysqlDB.Close()

			if err = mysqlDB.Ping(); err != nil {
				return fmt.Errorf("ping MySQL error: %w", err)
			}
			log.Println("[migrate] MySQL connected")

			// 确保 SQLite 目录存在
			if dir := filepath.Dir(sqlitePath); dir != "" {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("create sqlite dir error: %w", err)
				}
			}

			// 如果 SQLite 文件已存在，先删除
			if _, err := os.Stat(sqlitePath); err == nil {
				log.Printf("[migrate] existing SQLite file found, removing: %s", sqlitePath)
				_ = os.Remove(sqlitePath)
			}

			// 连接 SQLite
			sqliteDSN := fmt.Sprintf("file:%s?_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)", sqlitePath)
			sqliteDB, err := sql.Open("sqlite", sqliteDSN)
			if err != nil {
				return fmt.Errorf("open SQLite error: %w", err)
			}
			defer sqliteDB.Close()
			log.Println("[migrate] SQLite connected")

			// 导入 SQLite schema
			schemaPath := testutil.Schema()
			log.Printf("[migrate] importing schema: %s", schemaPath)
			if err := importSchema(sqliteDB, schemaPath); err != nil {
				return fmt.Errorf("import schema error: %w", err)
			}
			log.Println("[migrate] schema imported")

			// 逐表迁移数据
			var totalRows int64
			for _, table := range migrateTables {
				rows, err := migrateTable(ctx, mysqlDB, sqliteDB, table.name, table.columns)
				if err != nil {
					return fmt.Errorf("migrate table %s error: %w", table.name, err)
				}
				totalRows += rows
				log.Printf("[migrate] table %s: %d rows migrated", table.name, rows)
			}

			log.Printf("[migrate] completed! total rows: %d", totalRows)

			// 验证数据一致性
			if err := verifyMigration(ctx, mysqlDB, sqliteDB); err != nil {
				return fmt.Errorf("verification error: %w", err)
			}
			log.Println("[migrate] verification passed!")

			return nil
		},
	}
}

// importSchema 导入 SQLite 建表语句
func importSchema(db *sql.DB, schemaPath string) error {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	// 按分号分割并执行每条语句
	statements := splitSQLStatements(string(content))
	for _, stmt := range statements {
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("exec statement error: %w\nSQL: %s", err, stmt)
		}
	}
	return nil
}

// splitSQLStatements 将 SQL 文件按分号分割为独立语句（处理触发器的 END; ）
func splitSQLStatements(content string) []string {
	var statements []string
	var current strings.Builder
	inTrigger := false

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 跳过注释行
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}

		current.WriteString(line + "\n")

		// 检测触发器开始/结束
		upperLine := strings.ToUpper(trimmed)
		if strings.Contains(upperLine, "CREATE TRIGGER") {
			inTrigger = true
		}
		if inTrigger && strings.TrimSpace(trimmed) == "END;" {
			statements = append(statements, current.String())
			current.Reset()
			inTrigger = false
		} else if !inTrigger && strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if strings.TrimSpace(current.String()) != "" {
		statements = append(statements, current.String())
	}
	return statements
}

// migrateTable 迁移单个表的数据
func migrateTable(ctx context.Context, mysqlDB, sqliteDB *sql.DB, tableName string, columns []string) (int64, error) {
	// 从 MySQL 读取数据
	colList := "`" + strings.Join(columns, "`,`") + "`"
	query := fmt.Sprintf("SELECT %s FROM %s", colList, tableName)
	rows, err := mysqlDB.QueryContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("query %s error: %w", tableName, err)
	}
	defer rows.Close()

	// 准备 SQLite INSERT 语句
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName, colList, strings.Join(placeholders, ","))

	tx, err := sqliteDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx error: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return 0, fmt.Errorf("prepare insert error: %w", err)
	}
	defer stmt.Close()

	var count int64
	for rows.Next() {
		// 读取所有列
		scanArgs := make([]any, len(columns))
		for i := range scanArgs {
			scanArgs[i] = new(any)
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return count, fmt.Errorf("scan row error: %w", err)
		}

		// 转换参数值
		args := make([]any, len(columns))
		for i, arg := range scanArgs {
			val := *(arg.(*any))
			args[i] = convertValue(val)
		}

		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return count, fmt.Errorf("insert row error: %w", err)
		}
		count++
	}

	if err := rows.Err(); err != nil {
		return count, fmt.Errorf("rows error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return count, fmt.Errorf("commit error: %w", err)
	}

	// 更新 SQLite 的 autoincrement 序列值
	if count > 0 {
		updateAutoIncrement(ctx, sqliteDB, tableName, columns[0])
	}

	return count, nil
}

// convertValue 转换 MySQL 返回的值为 SQLite 兼容格式
func convertValue(v any) any {
	if v == nil {
		return nil
	}
	// []byte（MySQL JSON/BLOB 列）转为 string
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return v
}

// updateAutoIncrement 更新 SQLite 的 autoincrement 序列
func updateAutoIncrement(ctx context.Context, db *sql.DB, tableName, idCol string) {
	var maxID int
	err := db.QueryRowContext(ctx, fmt.Sprintf("SELECT MAX(%s) FROM %s", idCol, tableName)).Scan(&maxID)
	if err != nil || maxID == 0 {
		return
	}
	// 更新 sqlite_sequence 表
	_, _ = db.ExecContext(ctx, "UPDATE sqlite_sequence SET seq = ? WHERE name = ?", maxID, tableName)
}

// verifyMigration 验证迁移前后各表行数一致
func verifyMigration(ctx context.Context, mysqlDB, sqliteDB *sql.DB) error {
	for _, table := range migrateTables {
		var mysqlCount, sqliteCount int

		err := mysqlDB.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table.name)).Scan(&mysqlCount)
		if err != nil {
			return fmt.Errorf("count MySQL %s error: %w", table.name, err)
		}

		err = sqliteDB.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table.name)).Scan(&sqliteCount)
		if err != nil {
			return fmt.Errorf("count SQLite %s error: %w", table.name, err)
		}

		if mysqlCount != sqliteCount {
			return fmt.Errorf("row count mismatch for %s: MySQL=%d, SQLite=%d", table.name, mysqlCount, sqliteCount)
		}
	}
	return nil
}
