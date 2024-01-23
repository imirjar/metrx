package server

import (
	"testing"
)

// import (
// 	"testing"

// 	"github.com/imirjar/metrx/internal/service"
// )

func TestServerUpdateGauge(t *testing.T) {
	server := NewServerService()
	tests := []struct { // добавляем слайс тестов
		name     string
		mName    string
		mValue   float64
		expected error
	}{
		{
			name:     "Gauge",
			mName:    "SomeGauge",
			mValue:   123,
			expected: nil,
		},

		{
			name:     "Gauge without name",
			mName:    "",
			mValue:   11,
			expected: errMetricNameIncorrect, // must be name Error
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := server.UpdateGauge(test.mName, test.mValue)
			if err != test.expected {
				t.Errorf("%s, Error: %s", test.name, err)
			}
		})
	}
}

func TestServerUpdateCounter(t *testing.T) {
	server := NewServerService()
	tests := []struct { // добавляем слайс тестов
		name     string
		mName    string
		mValue   int64
		expected error
	}{
		{
			name:     "Gauge",
			mName:    "SomeGauge",
			mValue:   123,
			expected: nil,
		},

		{
			name:     "Gauge without name",
			mName:    "",
			mValue:   11,
			expected: errMetricNameIncorrect, // must be name Error
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := server.UpdateCounter(test.mName, test.mValue)
			if err != test.expected {
				t.Errorf("%s, Error: %s", test.name, err)
			}
		})
	}
}

func TestServerViewGauge(t *testing.T) {
	server := NewServerService()
	tests := []struct { // добавляем слайс тестов
		name          string
		mName         string
		mValue        float64
		expectedValue float64
		expectedErr   error
	}{
		{
			name:          "Gauge",
			mName:         "SomeGauge",
			mValue:        123,
			expectedValue: 123,
			expectedErr:   nil,
		},

		{
			name:          "Gauge without name",
			mName:         "",
			mValue:        11,
			expectedValue: 11,
			expectedErr:   errMetricNameIncorrect, // must be name Error
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.mName == "" {
				test.expectedValue = 0
				gauge, err := server.ViewGaugeByName(test.mName)
				if err != test.expectedErr || gauge != test.expectedValue {
					t.Errorf("\nValue: %f Expected: %f \nError: %s Expected: %s ", test.mValue, test.expectedValue, err, test.expectedErr)
				}
			} else {
				err := server.UpdateGauge(test.mName, test.mValue)
				if err != nil {
					t.Errorf("%s, Error: %s", test.name, err)
				}
				gauge, err := server.ViewGaugeByName(test.mName)
				if err != test.expectedErr || gauge != test.expectedValue {
					t.Errorf("\nValue: %f Expected: %f \nError: %s Expected: %s ", test.mValue, test.expectedValue, err, test.expectedErr)
				}
			}

		})
	}
}

func TestViewUpdateCounter(t *testing.T) {
	server := NewServerService()
	tests := []struct { // добавляем слайс тестов
		name     string
		mName    string
		mValue   int64
		expected error
	}{
		{
			name:     "Gauge",
			mName:    "SomeGauge",
			mValue:   123,
			expected: nil,
		},

		{
			name:     "Gauge without name",
			mName:    "",
			mValue:   11,
			expected: errMetricNameIncorrect, // must be name Error
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := server.UpdateCounter(test.mName, test.mValue)
			if err != test.expected {
				t.Errorf("%s, Error: %s", test.name, err)
			}
		})
	}
}
