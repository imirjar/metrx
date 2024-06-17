package models

import (
	"fmt"
	"strconv"
)

// Metrcis structure is used for containing system metric value, its name and type
type Metrics struct {
	ID    string   `json:"id" schema:"memberId"` // имя метрики
	MType string   `json:"type"`                 // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"`      // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"`      // значение метрики в случае передачи gauge
}

// Getter
func (m *Metrics) GetVal() (string, error) {
	switch m.MType {
	case "gauge":
		return fmt.Sprint(*m.Value), nil
	case "counter":
		return fmt.Sprint(*m.Delta), nil
	default:
		return "", errMetricTypeError
	}
}

// Setter
func (m *Metrics) SetVal(strVal string) error {
	switch m.MType {
	case "gauge":
		value, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return errTypeAssertionError
		}
		m.Value = &value
		return nil
	case "counter":
		delta, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return errTypeAssertionError
		}
		m.Delta = &delta
		return nil
	default:
		return errMetricTypeError
	}
}
