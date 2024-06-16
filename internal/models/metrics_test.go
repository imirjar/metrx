package models

import (
	"testing"
)

func TestMetrics_GetVal(t *testing.T) {
	tests := []struct {
		name    string
		metric  Metrics
		want    string
		wantErr bool
	}{
		{
			name: "gauge value",
			metric: Metrics{
				ID:    "testGauge",
				MType: "gauge",
				Value: float64Pointer(3.14),
			},
			want:    "3.14",
			wantErr: false,
		},
		{
			name: "counter value",
			metric: Metrics{
				ID:    "testCounter",
				MType: "counter",
				Delta: int64Pointer(42),
			},
			want:    "42",
			wantErr: false,
		},
		{
			name: "invalid type",
			metric: Metrics{
				ID:    "testInvalid",
				MType: "invalid",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.metric.GetVal()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_SetVal(t *testing.T) {
	tests := []struct {
		name    string
		metric  Metrics
		val     string
		wantErr bool
	}{
		{
			name: "set gauge value",
			metric: Metrics{
				ID:    "testGauge",
				MType: "gauge",
			},
			val:     "3.14",
			wantErr: false,
		},
		{
			name: "set counter value",
			metric: Metrics{
				ID:    "testCounter",
				MType: "counter",
			},
			val:     "42",
			wantErr: false,
		},
		{
			name: "set invalid gauge value",
			metric: Metrics{
				ID:    "testGauge",
				MType: "gauge",
			},
			val:     "invalid",
			wantErr: true,
		},
		{
			name: "set invalid counter value",
			metric: Metrics{
				ID:    "testCounter",
				MType: "counter",
			},
			val:     "invalid",
			wantErr: true,
		},
		{
			name: "invalid type",
			metric: Metrics{
				ID:    "testInvalid",
				MType: "invalid",
			},
			val:     "123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.metric.SetVal(tt.val); (err != nil) != tt.wantErr {
				t.Errorf("SetVal() error = %v, wantErr %v", err, tt.wantErr)
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
