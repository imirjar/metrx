package models

import (
	"math/rand"
	"reflect"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) SetMemValue(value reflect.Value, stat, mType string) error {

	m.ID = stat
	m.MType = mType
	if value.CanFloat() {
		val := value.Float()
		m.Value = &val
		// a.Storage.AddGauge(ms, value.Float())
	} else if value.CanUint() {
		val := float64(value.Uint())
		m.Value = &val
		// a.Storage.AddGauge(ms, float64(value.Uint()))

	}
	return nil
}

func (m *Metrics) SetRandomValue() {
	randV := rand.Float64()
	m.ID = "RandomValue"
	m.MType = "gauge"
	m.Value = &randV
}
