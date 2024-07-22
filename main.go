package main

import (
	"context"
	"orders-api/application"
	"os"
	"os/signal"
)

func main() {
	app := application.NewApp()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := app.Start(ctx)

	if err != nil {
		panic(err)
	}
}
