package usecase

import (
	"fmt"
	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/ws"
	"time"
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
	InsertAction(ActionType, string) error
	GetDevices() ([]Device, error)
	GetTelemetryByPeriod(string, time.Time, time.Time) ([]Telemetry, error)
	InsertTelemetry(Device, int32) error
}

type dataSender interface {
	Send(v ws.DataForSend)
}

type ProcessorDeps struct {
	TelemetryRepository telemetryRepository
	DataSender          dataSender
}

type Processor struct {
	telemetryRepository telemetryRepository
	dataSender          dataSender
}

func NewProcessor(d *ProcessorDeps) *Processor {
	return &Processor{
		telemetryRepository: d.TelemetryRepository,
		dataSender:          d.DataSender,
	}
}

func (p *Processor) Process(wet, temperature int32) error {
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

	return nil
}
