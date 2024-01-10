package service

import (
	"testing"
)

func TestGauge(t *testing.T) {
	service := New()
	tests := []struct { // добавляем слайс тестов
		mtype    string
		name     string
		value    map[string]string
		expected error
	}{
		{
			name: "Gauge",
			value: map[string]string{
				"mtype": "gauge",
				"name":  "gaugeMetric",
				"value": "10",
			},
			expected: nil,
		},
		{
			name: "Counter",
			value: map[string]string{
				"mtype": "counter",
				"name":  "counterMetric",
				"value": "10",
			},
			expected: nil,
		},
		{
			name: "Unknown metric type",
			value: map[string]string{
				"mtype": "unknown",
				"name":  "unknownMetric",
				"value": "10",
			},
			expected: serviceError,
		},
		{
			name: "Counter with incorrect value",
			value: map[string]string{
				"mtype": "counter",
				"name":  "counterMetric",
				"value": "counterValue",
			},
			expected: convertationError,
		},
		{
			name: "Gauge with incorrect value",
			value: map[string]string{
				"mtype": "gauge",
				"name":  "gaugeMetric",
				"value": "gaugeValue",
			},
			expected: convertationError,
		},

		{
			name: "Gauge without name",
			value: map[string]string{
				"mtype": "gauge",
				"name":  "",
				"value": "10",
			},
			expected: metricNameIncorrect, // must be name Error
		},
		{
			name: "Counter without name",
			value: map[string]string{
				"mtype": "counter",
				"name":  "",
				"value": "10",
			},
			expected: metricNameIncorrect, // must be name Error
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.UpdateMetric(test.value["mtype"], test.value["name"], test.value["value"])
			if err != test.expected {
				t.Errorf("%s, Error: %s", test.name, err)
			}
		})
	}
}