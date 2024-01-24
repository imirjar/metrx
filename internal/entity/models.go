package entity

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) SendJSONToPath(path string) error {
	client := &http.Client{
		Timeout: time.Second * 1,
	}

	mm, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := client.Post(path, "application/json", bytes.NewBuffer(mm))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (m *Metrics) SetRandomValue() {
	randValue := rand.Float64()
	m.Value = &randValue
}
