package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/app"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/config/configuration"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/config/credentials"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf, err := configuration.New()
	if err != nil {
		slog.Error("cannot create config", "error", err)
		return
	}

	cred, err := credentials.New()
	if err != nil {
		slog.Error("cannot create cred", "error", err)
		return
	}

	application, err := app.New(conf, cred)
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
