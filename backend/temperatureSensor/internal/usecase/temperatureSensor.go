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
	return &TemperatureSensor{}
}

func (t *TemperatureSensor) GenerateData() int32 {
	if t.isFreezing {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		val := int32(r.Intn(2) + 1)
		t.currentTemp -= val
		return t.currentTemp
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	val := int32(r.Intn(26) + 10)
	t.currentTemp = val

	return val
}

func (t *TemperatureSensor) ChangeMode() {
	t.isFreezing = !t.isFreezing
}
