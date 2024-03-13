package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) GetVal() (string, error) {
	switch m.MType {
	case "gauge":
		value := fmt.Sprint(*m.Value)
		return value, nil
	case "counter":
		delta := fmt.Sprint(*m.Delta)
		return delta, nil
	default:
		return "", fmt.Errorf("error incorrect metric type")
	}
}

func (m *Metrics) SetVal(strVal string) error {
	switch m.MType {
	case "gauge":
		value, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return fmt.Errorf("error metric isn't float64")
		}
		m.Value = &value
		return nil
	case "counter":
		delta, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return fmt.Errorf("error metric isn't float64")
		}
		m.Delta = &delta
		return nil
	default:
		return fmt.Errorf("error incorrect metric type")
	}
}

func (m *Metrics) Marshal() ([]byte, error) {
	if m.MType == "gauge" {
		nnVal := &struct {
			ID    string `json:"id"`
			MType string `json:"type"`
			Value string `json:"value,omitempty"`
		}{
			ID:    m.ID,
			MType: m.MType,
			Value: fmt.Sprint(*m.Value),
		}
		return json.Marshal(nnVal)
	} else {
		return json.Marshal(&m)
	}
}
