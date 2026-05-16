package usecase

import (
	"math/rand"
	"time"
)

type WetSensorDeps struct {
	DelayDataGen time.Duration
}

type WetSensor struct {
	delayDataGen time.Duration
}

func NewTemperatureSensor(d *WetSensorDeps) *WetSensor {
	return &WetSensor{
		delayDataGen: d.DelayDataGen,
	}
}

func (t *WetSensor) GenerateData() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int32(r.Intn(50) + 30)
}
