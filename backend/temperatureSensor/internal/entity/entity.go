package entity

import (
	"math/rand"
	"time"
)

type TemperatureSensor struct {
	Degrees int
}

func (s *TemperatureSensor) GenerateTemperature() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Degrees = r.Intn(26) + 10
}
