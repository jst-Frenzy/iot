package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jst-Frenzy/iot/backend/dataIntegrator/internal/infra/ws"
)

const (
	minTemperature = 18
	maxTemperature = 28

	minWet = 40
	maxWet = 70
)

type SourceType string

const (
	SourceTypeUser   SourceType = "user"
	SourceTypeSystem SourceType = "controller"
)

type ActionType string

const (
	ActionTypeOff ActionType = "off"
	ActionTypeOn  ActionType = "on"
)

type Device string

const (
	DeviceTemperatureSensor Device = "temperature_sensor"
	DeviceWetSensor         Device = "wet_sensor"
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
	InsertAction(ActionType, Device, SourceType) error
	GetDevices() ([]Device, error)
	GetTelemetryByPeriod(string, time.Time, time.Time) ([]Telemetry, error)
	InsertTelemetry(Device, int32) error
	SaveUserMove(ActionType, Device) error
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
	fanEnabled          bool
	pumpEnabled         bool
	telemetryRepository telemetryRepository
	dataSender          dataSender
	controllerManager   controllerManager
	mutex               sync.RWMutex
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

	if err := p.processFan(ctx, temperature); err != nil {
		return err
	}

	if err := p.processPump(ctx, wet); err != nil {
		return err
	}

	return nil
}

func (p *Processor) GetDevices() ([]Device, error) {
	return p.telemetryRepository.GetDevices()
}

func (p *Processor) GetTelemetryByPeriod(deviceName string, from time.Time, to time.Time) ([]Telemetry, error) {
	return p.telemetryRepository.GetTelemetryByPeriod(deviceName, from, to)
}

func (p *Processor) processFan(
	ctx context.Context,
	temperature int32,
) error {
	switch {
	case temperature > maxTemperature:
		return p.enableFan(ctx, SourceTypeSystem)

	case temperature < minTemperature:
		return p.disableFan(ctx, SourceTypeSystem)
	}

	return nil
}

func (p *Processor) processPump(
	ctx context.Context,
	wet int32,
) error {
	switch {
	case wet < minWet:
		return p.enablePump(ctx, SourceTypeSystem)

	case wet > maxWet:
		return p.disablePump(ctx, SourceTypeSystem)
	}

	return nil
}

func (p *Processor) ChangeFanMode() error {
	p.mutex.RLock()
	enabled := p.fanEnabled
	p.mutex.RUnlock()

	if enabled {
		return p.disableFan(context.Background(), SourceTypeUser)
	}

	return p.enableFan(context.Background(), SourceTypeUser)
}

func (p *Processor) ChangePumpMode() error {
	p.mutex.RLock()
	enabled := p.pumpEnabled
	p.mutex.RUnlock()

	if enabled {
		return p.disablePump(context.Background(), SourceTypeUser)
	}

	return p.enablePump(context.Background(), SourceTypeUser)
}

func (p *Processor) enableFan(ctx context.Context, source SourceType) error {
	p.mutex.Lock()

	if p.fanEnabled {
		p.mutex.Unlock()
		return nil
	}

	p.fanEnabled = true
	p.mutex.Unlock()

	if err := p.controllerManager.OnFan(ctx); err != nil {
		p.mutex.Lock()
		p.fanEnabled = false
		p.mutex.Unlock()

		return fmt.Errorf("cannot turn on fan: %w", err)
	}

	if err := p.telemetryRepository.InsertAction(
		ActionTypeOn,
		DeviceFan,
		source,
	); err != nil {
		return fmt.Errorf(
			"cannot insert fan on action: %w",
			err,
		)
	}

	if err := p.telemetryRepository.SaveUserMove(ActionTypeOn, DeviceFan); err != nil {
		slog.Error("cannot save user action", "error", err.Error())
	}

	return nil
}

func (p *Processor) disableFan(ctx context.Context, source SourceType) error {
	p.mutex.Lock()

	if !p.fanEnabled {
		p.mutex.Unlock()
		return nil
	}

	p.fanEnabled = false
	p.mutex.Unlock()

	if err := p.controllerManager.OffFan(ctx); err != nil {
		p.mutex.Lock()
		p.fanEnabled = true
		p.mutex.Unlock()

		return fmt.Errorf("cannot turn off fan: %w", err)
	}

	if err := p.telemetryRepository.InsertAction(
		ActionTypeOff,
		DeviceFan,
		source,
	); err != nil {
		return fmt.Errorf(
			"cannot insert fan off action: %w",
			err,
		)
	}

	if err := p.telemetryRepository.SaveUserMove(ActionTypeOff, DeviceFan); err != nil {
		slog.Error("cannot save user action", "error", err.Error())
	}

	return nil
}

func (p *Processor) enablePump(ctx context.Context, source SourceType) error {
	p.mutex.Lock()

	if p.pumpEnabled {
		p.mutex.Unlock()
		return nil
	}

	p.pumpEnabled = true
	p.mutex.Unlock()

	if err := p.controllerManager.OnPump(ctx); err != nil {
		p.mutex.Lock()
		p.pumpEnabled = false
		p.mutex.Unlock()

		return fmt.Errorf("cannot turn on pump: %w", err)
	}

	if err := p.telemetryRepository.InsertAction(
		ActionTypeOn,
		DevicePump,
		source,
	); err != nil {
		return fmt.Errorf(
			"cannot insert pump on action: %w",
			err,
		)
	}

	if err := p.telemetryRepository.SaveUserMove(ActionTypeOn, DevicePump); err != nil {
		slog.Error("cannot save user action", "error", err.Error())
	}

	return nil
}

func (p *Processor) disablePump(ctx context.Context, source SourceType) error {
	p.mutex.Lock()

	if !p.pumpEnabled {
		p.mutex.Unlock()
		return nil
	}

	p.pumpEnabled = false
	p.mutex.Unlock()

	if err := p.controllerManager.OffPump(ctx); err != nil {
		p.mutex.Lock()
		p.pumpEnabled = true
		p.mutex.Unlock()

		return fmt.Errorf("cannot turn off pump: %w", err)
	}

	if err := p.telemetryRepository.InsertAction(
		ActionTypeOff,
		DevicePump,
		source,
	); err != nil {
		return fmt.Errorf(
			"cannot insert pump off action: %w",
			err,
		)
	}

	if err := p.telemetryRepository.SaveUserMove(ActionTypeOff, DevicePump); err != nil {
		slog.Error("cannot save user action", "error", err.Error())
	}

	return nil
}
