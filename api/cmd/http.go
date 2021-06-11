package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goapt/gee"
	"github.com/urfave/cli"

	"app/router"
)

var HTTPCmd = cli.Command{
	Name:  "http",
	Usage: "http command eg: ./app http --addr=:8081",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Usage: "http listen ip:port",
		},
	},
	Action: func(ctx *cli.Context) error {
		if !ctx.IsSet("addr") {
			_ = ctx.Set("addr", ":8081")
		}

		serv := gee.Default()
		srv := &http.Server{
			Addr:    ctx.String("addr"),
			Handler: serv,
		}

		// router
		router.Route(serv)

		gee.RegisterShutDown(func(sig os.Signal) {
			ctxw, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Close()
			if err := srv.Shutdown(ctxw); err != nil {
				log.Fatal("HTTP Server Shutdown:", err)
			}
			log.Println("HTTP Server exiting")
		})

		log.Println("[HTTP] Server listen:" + ctx.String("addr"))
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP listen: %s\n", err)
		}

		return nil
	},
}

func init() {
	register(HTTPCmd)
}
