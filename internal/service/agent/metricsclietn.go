package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/imirjar/metrx/internal/models"
)

func (m *MetricsClient) POSTMetric(metric *models.Metrics) {
	mm, err := json.Marshal(metric)
	var buf bytes.Buffer

	if err != nil {
		panic(err)
	}

	gz := gzip.NewWriter(&buf)
	gz.Write(mm)
	gz.Close()

	req, err := http.NewRequest(http.MethodPost, m.Path, &buf)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := m.Client.Do(req)
	if err != nil {
		fmt.Printf(err.Error())
	}
	resp.Body.Close()
}
