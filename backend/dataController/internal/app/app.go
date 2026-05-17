package app

import (
	"context"
	"fmt"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/config/configuration"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/config/credentials"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/infra/grpc"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/infra/grpc/services"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/dataintegrator"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/fan"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/pump"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/temperatureSensor"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/wetSensor"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/usecase"
	"log/slog"
	"sync"
)

type App struct {
	collector  *usecase.Collector
	grpcServer *grpc.Server
}

func New(conf *configuration.Config, cred *credentials.Credentials) (*App, error) {
	dataIntegratorSender, err := dataintegrator.New(&dataintegrator.Config{
		Address:        cred.DataIntegrator.Address,
		MaxMessageSize: cred.DataIntegrator.MaxMessageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create data integrator sender: %w", err)
	}

	temperatureSender, err := temperatureSensor.New(&temperatureSensor.Config{
		Address:        cred.TemperatureSensor.Address,
		MaxMessageSize: cred.TemperatureSensor.MaxMessageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create temperature sender: %w", err)
	}

	wetSender, err := wetSensor.New(&wetSensor.Config{
		Address:        cred.WetSensor.Address,
		MaxMessageSize: cred.WetSensor.MaxMessageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create wet sender: %w", err)
	}

	pumpDevice, err := pump.New(&pump.Config{
		Address:        cred.Pump.Address,
		MaxMessageSize: cred.Pump.MaxMessageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create pump: %w", err)
	}

	fanDevice, err := fan.New(&fan.Config{
		Address:        cred.Fan.Address,
		MaxMessageSize: cred.Fan.MaxMessageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create fan: %w", err)
	}

	collector := usecase.NewCollector(&usecase.CollectorDeps{
		Delay:             conf.DelayDataGen,
		TemperatureSender: temperatureSender,
		WetSender:         wetSender,
		DataSender:        dataIntegratorSender,
		Pump:              pumpDevice,
		Fan:               fanDevice,
	})

	grpcServer, err := grpc.New(&grpc.Deps{
		Conf:   conf.GRPCServer,
		Logger: slog.With("server", "grpcServer"),
		Services: []grpc.Service{
			services.NewController(&services.ControllerDeps{
				Manager: collector,
			}),
		},
	})

	return &App{
		collector:  collector,
		grpcServer: grpcServer,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	a.collector.Start(ctx)

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
