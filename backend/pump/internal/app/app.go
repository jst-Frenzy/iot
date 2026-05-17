package app

import (
	"context"
	"fmt"
	"github.com/jst-Frenzy/iot/backend/pump/internal/config/configuration"
	"github.com/jst-Frenzy/iot/backend/pump/internal/infra/grpc"
	"github.com/jst-Frenzy/iot/backend/pump/internal/infra/grpc/services"
	"log/slog"
	"sync"
)

type App struct {
	grpcServer *grpc.Server
}

func New(ctx context.Context, conf *configuration.Config) (*App, error) {
	grpcServer, err := grpc.New(&grpc.Deps{
		Conf:   conf.GRPCServer,
		Logger: slog.With("server", "grpcServer"),
		Services: []grpc.Service{
			services.New(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create grpcServer: %w", err)
	}

	return &App{
		grpcServer: grpcServer,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	errChan := make(chan error, 1)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.grpcServer.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		return err
	}

	return nil
}
