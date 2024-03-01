package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
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

func (m *Metrics) MarshalGauge() ([]byte, error) {
	if m.MType == "gauge" {
		val := strconv.FormatFloat(float64(*m.Value), 'E', -1, 64)
		return json.Marshal(&struct {
			ID    string `json:"id"`
			MType string `json:"type"`
			Value string `json:"value,omitempty"`
		}{
			ID:    m.ID,
			MType: m.MType,
			Value: fmt.Sprint(val),
		})
	} else {
		return json.Marshal(&m)
	}
}
