package usecase

import (
	"math/rand"
	"time"
)

type WetSensor struct{}

func NewTemperatureSensor() *WetSensor {
	return &WetSensor{}
}

func (t *WetSensor) GenerateData() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int32(r.Intn(50) + 30)
}
