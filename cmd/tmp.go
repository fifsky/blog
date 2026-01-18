package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"app/config"

	"github.com/urfave/cli/v3"
)

func NewTmp(db *sql.DB, conf *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "tmp",
		Usage: "tmp command eg: ./app tmp",
		Action: func(ctx context.Context, cli *cli.Command) error {
			fmt.Println("this is tmp command")
			return nil
		},
	}
}
