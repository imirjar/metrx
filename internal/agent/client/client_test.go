package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/imirjar/metrx/internal/models"
)

func TestClient_POST(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		metrics    []models.Metrics
		wantStatus int
	}{
		{
			name:   "valid request without secret",
			secret: "",
			metrics: []models.Metrics{
				{ID: "TestMetric1", MType: "gauge", Value: float64Pointer(3.14)},
				{ID: "TestMetric2", MType: "counter", Delta: int64Pointer(42)},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "valid request with secret",
			secret: "mysecret",
			metrics: []models.Metrics{
				{ID: "TestMetric3", MType: "gauge", Value: float64Pointer(6.28)},
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый HTTP-сервер
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected method POST, got %v", r.Method)
				}

				// Проверяем заголовки
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type application/json, got %v", r.Header.Get("Content-Type"))
				}
				if r.Header.Get("Content-Encoding") != "gzip" {
					t.Errorf("expected Content-Encoding gzip, got %v", r.Header.Get("Content-Encoding"))
				}

				// Проверяем тело запроса
				var buf bytes.Buffer
				gr, err := gzip.NewReader(r.Body)
				if err != nil {
					t.Fatalf("failed to create gzip reader: %v", err)
				}
				buf.ReadFrom(gr)
				gr.Close()

				var receivedMetrics []models.Metrics
				if err := json.Unmarshal(buf.Bytes(), &receivedMetrics); err != nil {
					t.Fatalf("failed to unmarshal metrics: %v", err)
				}

				// Проверяем, что метрики совпадают
				if len(receivedMetrics) != len(tt.metrics) {
					t.Errorf("expected %d metrics, got %d", len(tt.metrics), len(receivedMetrics))
				}

				// Отправляем ответ
				w.WriteHeader(tt.wantStatus)
			}))
			defer ts.Close()

			client := New(tt.secret, "", ts.URL)

			// Создаем контекст с тайм-аутом
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			err := client.POST(ctx, tt.metrics)
			if err != nil {
				t.Fatalf("POST() error = %v", err)
			}
		})
	}
}

// Вспомогательные функции для указателей
func float64Pointer(v float64) *float64 {
	return &v
}

func int64Pointer(v int64) *int64 {
	return &v
}
