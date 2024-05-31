package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game-random-api/internal/api"
)

type keyType string

const ENV_KEY keyType = "mode"

func main() {
	ctx := context.Background()

	app_mode := os.Getenv("APP_MODE")
	if app_mode == "" {
		app_mode = "development"
	}

	appCtx := context.WithValue(ctx, ENV_KEY, app_mode)

	app := api.NewAPI(appCtx)

	ctxDone, cancel := context.WithCancel(ctx)
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer cancel()
		if err := app.Setup(ctx); err != nil {
			fmt.Println("Failed to setup application", err.Error())
			return
		}
		app.Start(ctx)
	}(appCtx, cancel)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 2)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		fmt.Println("Shutdown Server via external signal")
		app.Interrupt(ctxDone)
	case <-ctxDone.Done():
		fmt.Println("Shutdown server via internal request")
	}

	quitCtx, done := context.WithTimeout(ctx, 10*time.Second)
	go func() {
		defer done()
		app.Stop(quitCtx)
		app.Deinit(quitCtx)
	}()
	// catching ctx.Done(). timeout of 5 seconds.
	<-quitCtx.Done()
}
