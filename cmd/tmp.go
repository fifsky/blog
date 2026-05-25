package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"app/config"

	"github.com/urfave/cli/v3"
)

func NewTmp(db *sql.DB, conf *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "tmp",
		Usage: "migrate reminds cron data",
		Action: func(ctx context.Context, cli *cli.Command) error {
			fmt.Println("Start migrating reminds table to cron...")

			// Query old data
			rows, err := db.QueryContext(ctx, "SELECT id, type, month, week, day, hour, minute, created_at FROM reminds")
			if err != nil {
				return fmt.Errorf("query reminds error: %w", err)
			}
			defer rows.Close()

			var updates []struct {
				id   int
				cron string
			}

			location, _ := time.LoadLocation("Asia/Shanghai")
			if location == nil {
				location = time.Local
			}

			for rows.Next() {
				var id, typ, month, week, day, hour, minute int
				var createdAt time.Time

				if err := rows.Scan(&id, &typ, &month, &week, &day, &hour, &minute, &createdAt); err != nil {
					return fmt.Errorf("scan error: %w", err)
				}

				var cronStr string
				switch typ {
				case 0:
					year := createdAt.In(location).Year()
					cronStr = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", year, month, day, hour, minute)
				case 1:
					cronStr = "* * * * *"
				case 2:
					cronStr = fmt.Sprintf("%d * * * *", minute)
				case 3:
					cronStr = fmt.Sprintf("%d %d * * *", minute, hour)
				case 4:
					cronWeek := week
					if cronWeek == 7 {
						cronWeek = 0
					}
					cronStr = fmt.Sprintf("%d %d * * %d", minute, hour, cronWeek)
				case 5:
					cronStr = fmt.Sprintf("%d %d %d * *", minute, hour, day)
				case 6:
					cronStr = fmt.Sprintf("%d %d %d %d *", minute, hour, day, month)
				}

				updates = append(updates, struct {
					id   int
					cron string
				}{id, cronStr})
			}

			if err := rows.Err(); err != nil {
				return err
			}

			// Update rows
			for _, u := range updates {
				_, err := db.ExecContext(ctx, "UPDATE reminds SET cron = ? WHERE id = ?", u.cron, u.id)
				if err != nil {
					return fmt.Errorf("update id %d error: %w", u.id, err)
				}
				fmt.Printf("Updated id %d with cron: %s\n", u.id, u.cron)
			}

			fmt.Println("Migration completed successfully!")
			return nil
		},
	}
}
