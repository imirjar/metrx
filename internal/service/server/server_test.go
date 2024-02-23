package server

// import (
// 	"testing"

// 	"github.com/imirjar/metrx/config"
// 	"github.com/imirjar/metrx/internal/models"
// )

// func TestServerGauge(t *testing.T) {
// 	cfg := config.Testcfg
// 	server := NewServerService(cfg)

// 	tests := []struct { // добавляем слайс тестов
// 		name          string
// 		mName         string
// 		metric        models.Metrics
// 		expectedValue float64
// 		expectedErr   error
// 	}{
// 		{
// 			name:  "Gauge",
// 			mName: "SomeGauge",

// 			expectedValue: 123,
// 			expectedErr:   nil,
// 		},
// 	}
// 	for _, test := range tests {
// 		var value float64 = 123
// 		metric := models.Metrics{
// 			ID:    "gauge",
// 			MType: "gauge",
// 			Value: &value,
// 		}
// 		t.Run(test.name, func(t *testing.T) {
// 			nMetric, err := server.Update(metric)
// 			if err != test.expectedErr {
// 				t.Errorf("%s, Error: %s", test.name, err)
// 			}
// 			gauge, err := server.View(metric)
// 			if err != test.expectedErr || gauge != nMetric {
// 				t.Errorf("\nValue: %f Expected: %f \nError: %s Expected: %s ", *metric.Value, test.expectedValue, err, test.expectedErr)
// 			}
// 		})
// 	}
// }
