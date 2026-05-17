package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type temperatureSender interface {
	GetTemperature(context.Context) (int32, error)
}

type wetSender interface {
	GetWet(context.Context) (int32, error)
}

type dataSender interface {
	SendData(context.Context, int32, int32) error
}

type CollectorDeps struct {
	Delay             time.Duration
	TemperatureSender temperatureSender
	WetSender         wetSender
	DataSender        dataSender
}

type Collector struct {
	delay             time.Duration
	temperatureSender temperatureSender
	wetSender         wetSender
	dataSender        dataSender
	logger            *slog.Logger
}

func NewCollector(d *CollectorDeps) *Collector {
	return &Collector{
		delay:             d.Delay,
		temperatureSender: d.TemperatureSender,
		wetSender:         d.WetSender,
		dataSender:        d.DataSender,
		logger:            slog.With("collector", "collect data"),
	}
}

func (c *Collector) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(c.delay)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("collector stopped")
				return
			case <-ticker.C:
			}

			wet, err := c.wetSender.GetWet(ctx)
			if err != nil {
				c.logger.Error("cant get wet", "error", err)
				continue
			}

			temperature, err := c.temperatureSender.GetTemperature(ctx)
			if err != nil {
				c.logger.Error("cant get temperature", "error", err)
				continue
			}

			fmt.Println("wet: ", wet, "temperature: ", temperature)

			err = c.dataSender.SendData(ctx, wet, temperature)
			if err != nil {
				c.logger.Error("cant send data to data integrator", "error", err)
			}

			c.logger.Info("successfully send data to data integrator")
		}
	}()
}
