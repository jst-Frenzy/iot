package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type temperatureSender interface {
	GetTemperature(context.Context) (int32, error)
	ChangeMode(ctx context.Context) error
}

type wetSender interface {
	GetWet(context.Context) (int32, error)
	ChangeMode(ctx context.Context) error
}

type dataSender interface {
	SendData(context.Context, int32, int32) error
}

type pump interface {
	On(context.Context) error
	Off(context.Context) error
}

type fan interface {
	On(context.Context) error
	Off(context.Context) error
}

type CollectorDeps struct {
	Delay             time.Duration
	TemperatureSender temperatureSender
	WetSender         wetSender
	DataSender        dataSender
	Pump              pump
	Fan               fan
}

type Collector struct {
	delay             time.Duration
	temperatureSender temperatureSender
	wetSender         wetSender
	dataSender        dataSender
	pump              pump
	fan               fan
	logger            *slog.Logger
}

func NewCollector(d *CollectorDeps) *Collector {
	return &Collector{
		delay:             d.Delay,
		temperatureSender: d.TemperatureSender,
		wetSender:         d.WetSender,
		dataSender:        d.DataSender,
		pump:              d.Pump,
		fan:               d.Fan,
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

func (c *Collector) OnPump(ctx context.Context) error {
	err := c.pump.On(ctx)
	if err != nil {
		return fmt.Errorf("cannot on pump: %w", err)
	}

	err = c.wetSender.ChangeMode(ctx)
	if err != nil {
		return fmt.Errorf("cannot change wet mode: %w", err)
	}

	return nil
}

func (c *Collector) OffPump(ctx context.Context) error {
	err := c.pump.Off(ctx)
	if err != nil {
		return fmt.Errorf("cannot off pump: %w", err)
	}

	err = c.wetSender.ChangeMode(ctx)
	if err != nil {
		return fmt.Errorf("cannot change wet mode: %w", err)
	}

	return nil
}

func (c *Collector) OnFan(ctx context.Context) error {
	err := c.fan.On(ctx)
	if err != nil {
		return fmt.Errorf("cannot on fan: %w", err)
	}

	err = c.temperatureSender.ChangeMode(ctx)
	if err != nil {
		return fmt.Errorf("cannot change fan mode: %w", err)
	}

	return nil
}

func (c *Collector) OffFan(ctx context.Context) error {
	err := c.fan.Off(ctx)
	if err != nil {
		return fmt.Errorf("cannot off fan: %w", err)
	}

	err = c.temperatureSender.ChangeMode(ctx)
	if err != nil {
		return fmt.Errorf("cannot change fan mode: %w", err)
	}

	return nil
}
