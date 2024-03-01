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

func stringValue(f *float64) *string {
	if f != nil {
		strValue := fmt.Sprint(*f)
		return &strValue
	}
	return nil
}

func (m *Metrics) MarshalJSON() ([]byte, error) {
	type Alias Metrics
	return json.Marshal(&struct {
		*Alias
		Value *string `json:"value,omitempty"`
	}{
		Alias: (*Alias)(m),
		Value: stringValue(m.Value),
	})
}
