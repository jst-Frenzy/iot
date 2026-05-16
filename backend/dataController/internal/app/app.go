package app

import (
	"context"
	"fmt"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/config/configuration"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/config/credentials"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/dataintegrator"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/temperatureSensor"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/integrations/grpc/wetSensor"
	"github.com/jst-Frenzy/iot/backend/dataController/internal/usecase"
)

type App struct {
	collector *usecase.Collector
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

	collector := usecase.NewCollector(&usecase.CollectorDeps{
		Delay:             conf.DelayDataGen,
		TemperatureSender: temperatureSender,
		WetSender:         wetSender,
		DataSender:        dataIntegratorSender,
	})

	return &App{
		collector: collector,
	}, nil
}

func (a *App) Start(ctx context.Context) {
	a.collector.Start(ctx)
}
