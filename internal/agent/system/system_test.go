package system

import (
	"context"
	"runtime"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	collector := NewSystem()

	metrics, err := collector.Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	if len(metrics) == 0 {
		t.Fatal("Expected non-zero number of metrics")
	}

	// Проверка некоторых метрик на наличие
	expectedMetrics := []string{"Alloc", "BuckHashSys", "Frees"}
	for _, name := range expectedMetrics {
		found := false
		for _, metric := range metrics {
			if metric.ID == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Metric %s not found in collected metrics", name)
		}
	}
}

func Test_getMemStat(t *testing.T) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	tests := []struct {
		name      string
		metric    string
		wantError bool
	}{
		{"Alloc", "Alloc", false},
		{"BuckHashSys", "BuckHashSys", false},
		{"Frees", "Frees", false},
		{"InvalidMetric", "InvalidMetric", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getMemStat(&ms, tt.metric)
			if (err != nil) != tt.wantError {
				t.Errorf("getMemStat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
