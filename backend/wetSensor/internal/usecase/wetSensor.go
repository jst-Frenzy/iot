package usecase

import (
	"math/rand"
	"time"
)

type WetSensor struct {
	isWetting  bool
	currentWet int32
}

func NewTemperatureSensor() *WetSensor {
	return &WetSensor{
		isWetting:  false,
		currentWet: 55,
	}
}

func (t *WetSensor) GenerateData() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	val := int32(r.Intn(3) + 1)

	if t.isWetting {
		t.currentWet += val
	}
	t.currentWet -= val

	return t.currentWet
}

func (t *WetSensor) ChangeMode() {
	t.isWetting = !t.isWetting
}
