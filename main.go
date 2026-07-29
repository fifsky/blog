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
	// 显式调用 Close，确保即使 Run 返回错误（log.Fatal 会 os.Exit 跳过 defer）也能释放资源
	err := app.Run(ctx)
	if cerr := app.Close(); cerr != nil {
		log.Println(cerr)
	}
	if err != nil {
		log.Fatal(err)
	}
}
