package usecase

import (
	"github.com/jst-Frenzy/iot/backend/temperatureSensor/internal/entity"
	"time"
)

type TemperatureSensorDeps struct {
	DelayDataGen time.Duration
}

type TemperatureSensor struct {
	device       entity.TemperatureSensor
	delayDataGen time.Duration
}

func NewTemperatureSensor(d *TemperatureSensorDeps) *TemperatureSensor {
	return &TemperatureSensor{
		device:       entity.TemperatureSensor{},
		delayDataGen: d.DelayDataGen,
	}
}

func (t *TemperatureSensor) Start() {
	go func() {
		time.Sleep(t.delayDataGen)
		t.device.GenerateTemperature()

	}()
}
