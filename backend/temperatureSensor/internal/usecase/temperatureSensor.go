package usecase

import (
	"math/rand"
	"time"
)

type TemperatureSensor struct {
}

func NewTemperatureSensor() *TemperatureSensor {
	return &TemperatureSensor{}
}

func (t *TemperatureSensor) GenerateData() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int32(r.Intn(26) + 10)
}
