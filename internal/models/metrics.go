package models

import (
	"encoding/json"
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
	if m.MType == "gauge" {
		return json.Marshal(&struct {
			ID    string `json:"id"`   // имя метрики
			MType string `json:"type"` // параметр, принимающий значение gauge или counter
			Value string `json:"value,omitempty"`
		}{
			ID:    m.ID,
			MType: m.MType,
			Value: fmt.Sprintf("%f", *m.Value),
		})
	} else {
		return json.Marshal(&m)
	}
}
