package models

import (
	"encoding/json"
	"fmt"
	"log"
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

func (m *Metrics) MarshalGauge() ([]byte, error) {
	if m.MType == "gauge" {
		// val := strconv.FormatFloat(*m.Value, 'f', -1, 64)
		log.Println("#######MarshalGauge+++++++++++", *m.Value)
		nnVal := &struct {
			ID    string `json:"id"`
			MType string `json:"type"`
			Value string `json:"value,omitempty"`
		}{
			ID:    m.ID,
			MType: m.MType,
			Value: fmt.Sprintln(*m.Value),
		}
		log.Println("#######AS_YOBANIY_RESULT+++++++++++", nnVal.Value)
		return json.Marshal(nnVal)
	} else {
		return json.Marshal(&m)
	}
}
