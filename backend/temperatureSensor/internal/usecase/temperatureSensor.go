package usecase

import (
	"math/rand"
	"time"
)

type TemperatureSensorDeps struct {
	DelayDataGen time.Duration
}

type TemperatureSensor struct {
	delayDataGen time.Duration
}

func NewTemperatureSensor(d *TemperatureSensorDeps) *TemperatureSensor {
	return &TemperatureSensor{
		delayDataGen: d.DelayDataGen,
	}
}

func (t *TemperatureSensor) GenerateData() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int32(r.Intn(26) + 10)
}
