package usecase

import (
	"math/rand"
	"time"
)

type TemperatureSensor struct {
	isFreezing  bool
	currentTemp int32
}

func NewTemperatureSensor() *TemperatureSensor {
	return &TemperatureSensor{
		isFreezing:  false,
		currentTemp: 22,
	}
}

func (t *TemperatureSensor) GenerateData() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	val := int32(r.Intn(2) + 1)

	if t.isFreezing {
		t.currentTemp -= val
	}
	t.currentTemp += val

	return t.currentTemp
}

func (t *TemperatureSensor) ChangeMode() {
	t.isFreezing = !t.isFreezing
}
