package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/imirjar/metrx/internal/models"
)

func (m *MetricsClient) POSTMetric(metric models.Metrics) error {

	mm, err := json.Marshal(metric)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(mm)
	gz.Close()

	path := fmt.Sprintf(m.Path + "/update/")
	req, err := http.NewRequest(http.MethodPost, path, &buf)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := m.Client.Do(req)

	if err != nil {
		log.Fatal(err)
		return err
	}

	resp.Body.Close()
	return err
}

func (m *MetricsClient) POSTMetrics(metric []models.Metrics) error {
	for me := range metric {
		mem := metric[me]
		if mem.MType == "gauge" {
			log.Println(mem.MType, mem.ID, *mem.Value)
		} else {
			log.Println(mem.MType, mem.ID, *mem.Delta)
		}

	}
	mm, err := json.Marshal(metric)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(mm)
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, m.Path+"/updates/", &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := m.Client.Do(req)

	if err != nil {
		return err
	}

	resp.Body.Close()
	return err
}
