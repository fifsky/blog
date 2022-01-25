package cmd

import (
	"log"

	"app/config"
	"github.com/urfave/cli/v2"

	"app/router"
)

type HttpCmd *cli.Command

func NewHttp(router router.Router, conf *config.Config) HttpCmd {
	return &cli.Command{
		Name:  "http",
		Usage: "http command eg: ./app http --addr=:8080",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Usage: "http listen ip:port",
			},
		},
		Action: func(ctx *cli.Context) error {
			if !ctx.IsSet("addr") {
				_ = ctx.Set("addr", ":8080")
			}
			log.Println("[Env] Run profile:" + conf.Env)
			router.Run(ctx.String("addr"))
			return nil
		},
	}
}
