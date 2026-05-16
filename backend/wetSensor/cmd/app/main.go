package main

import (
	"context"
	"github.com/jst-Frenzy/iot/backend/wetSensor/internal/app"
	"github.com/jst-Frenzy/iot/backend/wetSensor/internal/config/configuration"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := configuration.New()
	if err != nil {
		slog.Error("cannot create config", "error", err)
		return
	}

	application, err := app.New(ctx, conf)
	if err != nil {
		slog.Error("cannot create app", "error", err)
		return
	}

	err = application.Start(ctx)
	if err != nil {
		slog.Error("cannot create app", "error", err)
		return
	}

	slog.Info("started")
	<-ctx.Done()
}
