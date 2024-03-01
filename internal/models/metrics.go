package models

import (
	"fmt"
	"math/rand"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) SetRandomValue() {
	randV := rand.Float64()
	m.Value = &randV
}

func (m *Metrics) MarshalJSON() ([]byte, error) {
	switch m.MType {
	case "gauge":
		return []byte(fmt.Sprintf(`{"ID":%s,"MType":%s,"Value":%f}`, m.ID, m.MType, *m.Value)), nil
	case "counter":
		return []byte(fmt.Sprintf(`{"ID":%s,"MType":%s,"Delta":%d}`, m.ID, m.MType, *m.Delta)), nil
	default:
		return []byte{}, nil
	}
}
