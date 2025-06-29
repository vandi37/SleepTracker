package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vandi37/SleepTracker/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := app.New()
	app.Run(ctx)
}
