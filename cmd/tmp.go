package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"app/config"

	"github.com/urfave/cli/v3"
)

// NewTmp 临时命令模板，用于一次性数据修复和临时脚本等场景
func NewTmp(db *sql.DB, conf *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "tmp",
		Usage: "临时命令，按需修改逻辑后执行",
		Action: func(ctx context.Context, cli *cli.Command) error {
			fmt.Println("tmp command running...")
			// 在此处编写临时逻辑
			return nil
		},
	}
}
