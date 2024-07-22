package application

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func NewApp() *App {
	app := &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
	return app
}

func (app *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: app.router,
	}
	err := app.rdb.Ping(ctx).Err()

	if err != nil {
		return fmt.Errorf("failed to connect to redis %w", err)
	}
	defer func() {
		if err := app.rdb.Close(); err != nil {
			fmt.Println("failed to close redis connection")
		}
	}()

	fmt.Println("Redis server is running")

	ch := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server :$w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(timeout)
	}

	return nil
}
