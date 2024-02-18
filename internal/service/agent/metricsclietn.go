package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/imirjar/metrx/internal/models"
)

func (m *MetricsClient) POSTMetric(metric models.Metrics) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(metric)
	if err != nil {
		log.Fatal(err)
	}

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
