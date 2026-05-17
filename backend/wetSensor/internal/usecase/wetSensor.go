package usecase

import (
	"math/rand"
	"time"
)

type WetSensor struct {
	iWetting   bool
	currentWet int32
}

func NewTemperatureSensor() *WetSensor {
	return &WetSensor{}
}

func (t *WetSensor) GenerateData() int32 {
	if t.iWetting {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		val := int32(r.Intn(2) + 1)
		t.currentWet += val
		return t.currentWet
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	val := int32(r.Intn(50) + 30)
	t.currentWet = val

	return val
}

func (t *WetSensor) ChangeMode() {
	t.iWetting = !t.iWetting
}
