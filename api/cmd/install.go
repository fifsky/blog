package cmd

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"

	"app/config"
)

type InstallCmd *cli.Command

func NewInstallCmd() InstallCmd {
	return &cli.Command{
		Name:  "install",
		Usage: "install command eg: ./app install",
		Action: func(ctx *cli.Context) error {
			_, err := config.ImportDB()
			if err != nil {
				fmt.Println("Import DB Error:" + err.Error())
				log.Fatalf("import error %s", err)
			}
			fmt.Println("Database init success!")
			return nil
		},
	}
}
