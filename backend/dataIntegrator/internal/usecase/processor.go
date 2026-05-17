package usecase

import (
	"context"
	"fmt"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/ws"
	"time"
)

const (
	minTemperature = 18
	maxTemperature = 28

	minWet = 40
	maxWet = 70
)

type ActionType string

const (
	ActionTypeOff ActionType = "off"
	ActionTypeOn  ActionType = "on"
)

type Device string

const (
	DeviceTemperatureSensor Device = "temperature sensor"
	DeviceWetSensor         Device = "wet sensor"
	DeviceFan               Device = "fan"
	DevicePump              Device = "pump"
)

type Telemetry struct {
	ID         int       `gorm:"column:id"`
	DeviceName string    `gorm:"column:device_name"`
	Value      int       `gorm:"column:value"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

type telemetryRepository interface {
	InsertAction(ActionType, Device) error
	GetDevices() ([]Device, error)
	GetTelemetryByPeriod(string, time.Time, time.Time) ([]Telemetry, error)
	InsertTelemetry(Device, int32) error
}

type dataSender interface {
	Send(v ws.DataForSend)
}

type controllerManager interface {
	OffPump(ctx context.Context) error
	OnPump(ctx context.Context) error
	OffFan(ctx context.Context) error
	OnFan(ctx context.Context) error
}

type ProcessorDeps struct {
	TelemetryRepository telemetryRepository
	DataSender          dataSender
	ControllerManager   controllerManager
}

type Processor struct {
	isFreezing          bool
	isWetting           bool
	telemetryRepository telemetryRepository
	dataSender          dataSender
	controllerManager   controllerManager
}

func NewProcessor(d *ProcessorDeps) *Processor {
	return &Processor{
		telemetryRepository: d.TelemetryRepository,
		dataSender:          d.DataSender,
		controllerManager:   d.ControllerManager,
	}
}

func (p *Processor) Process(ctx context.Context, wet, temperature int32) error {
	err := p.telemetryRepository.InsertTelemetry(DeviceTemperatureSensor, temperature)
	if err != nil {
		return fmt.Errorf("cannot insert temperature: %w", err)
	}

	err = p.telemetryRepository.InsertTelemetry(DeviceWetSensor, wet)
	if err != nil {
		return fmt.Errorf("cannot insert wet: %w", err)
	}

	p.dataSender.Send(ws.DataForSend{
		Temperature: temperature,
		Wet:         wet,
	})

	if !p.isFreezing && temperature > maxTemperature {
		err := p.controllerManager.OnFan(ctx)
		if err != nil {
			return fmt.Errorf("cannot turn on fan: %w", err)
		}

		err = p.telemetryRepository.InsertAction(ActionTypeOn, DeviceFan)
		if err != nil {
			return fmt.Errorf("cannot insert action: %w", err)
		}

		p.isFreezing = true
	}

	if p.isFreezing && temperature < minTemperature {
		err := p.controllerManager.OffFan(ctx)
		if err != nil {
			return fmt.Errorf("cannot turn off fan: %w", err)
		}

		err = p.telemetryRepository.InsertAction(ActionTypeOff, DeviceFan)
		if err != nil {
			return fmt.Errorf("cannot insert action: %w", err)
		}

		p.isFreezing = false
	}

	if !p.isWetting && wet < minWet {
		err := p.controllerManager.OnPump(ctx)
		if err != nil {
			return fmt.Errorf("cannot turn on pump: %w", err)
		}

		err = p.telemetryRepository.InsertAction(ActionTypeOn, DevicePump)
		if err != nil {
			return fmt.Errorf("cannot insert action: %w", err)
		}

		p.isWetting = true
	}

	if p.isWetting && wet > maxWet {
		err := p.controllerManager.OffPump(ctx)
		if err != nil {
			return fmt.Errorf("cannot turn off pump: %w", err)
		}

		err = p.telemetryRepository.InsertAction(ActionTypeOff, DeviceFan)
		if err != nil {
			return fmt.Errorf("cannot insert action: %w", err)
		}

		p.isWetting = false
	}

	return nil
}
