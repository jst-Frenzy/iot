package usecase

import (
	"context"
	"log/slog"
	"time"
)

type temperatureSender interface {
	GetTemperature(context.Context) (int32, error)
}

type wetSender interface {
	GetWet(context.Context) (int32, error)
}

type CollectorDeps struct {
	Delay             time.Duration
	TemperatureSender temperatureSender
	WetSender         wetSender
}

type Collector struct {
	delay             time.Duration
	temperatureSender temperatureSender
	wetSender         wetSender
	logger            *slog.Logger
}

func NewCollector(d *CollectorDeps) *Collector {
	return &Collector{
		delay:             d.Delay,
		temperatureSender: d.TemperatureSender,
		wetSender:         d.WetSender,
		logger:            slog.With("collector", "collect data"),
	}
}

func (c *Collector) Start(ctx context.Context) {
	go func() {
		time.Sleep(c.delay)
		wet, err := c.wetSender.GetWet(ctx)
		if err != nil {
			c.logger.Error("cant get wet", "error", err)
		}

		temperature, err := c.temperatureSender.GetTemperature(ctx)
		if err != nil {
			c.logger.Error("cant get temperature", "error", err)
		}

	}()
}
