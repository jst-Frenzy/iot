package app

import (
	"context"
	"fmt"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/config/configuration"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/config/credentials"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/data_base"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/grpc"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/grpc/services"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/postgres/repositories"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/ws"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/usecase"
	"log/slog"
	"sync"
)

type App struct {
	grpcServer *grpc.Server
}

func New(conf *configuration.Config, cred *credentials.Credentials) (*App, error) {
	postgres, err := data_base.InitPostgres(cred)
	if err != nil {
		return nil, fmt.Errorf("cannot init db")
	}

	repo := repositories.NewPostgres(&repositories.PostgresDeps{
		DB: postgres,
	})

	dataSender := ws.New(conf.WsServer.Address)

	processor := usecase.NewProcessor(&usecase.ProcessorDeps{
		TelemetryRepository: repo,
		DataSender:          dataSender,
	})

	grpcServer, err := grpc.New(&grpc.Deps{
		Conf:   conf.GRPCServer,
		Logger: slog.With("server", "grpcServer"),
		Services: []grpc.Service{
			services.NewAccepter(&services.AccepterDeps{
				Processor: processor,
			}),
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
