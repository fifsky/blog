package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"app/cmd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer stop()
	app := cmd.NewCommand()
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
