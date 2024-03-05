package entity

import (
	"bytes"
	"compress/gzip"
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

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(mm)
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, path, &buf)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (m *Metrics) SetRandomValue() {
	randValue := rand.Float64()
	m.Value = &randValue
}
