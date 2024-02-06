package server

import (
	"testing"

	"github.com/imirjar/metrx/config"
)

func TestServerGauge(t *testing.T) {
	cfg := config.Testcfg
	server := NewServerService(cfg)
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := server.UpdateGauge(test.mName, test.mValue)
			if err != test.expectedErr {
				t.Errorf("%s, Error: %s", test.name, err)
			}
			gauge, err := server.ViewGaugeByName(test.mName)
			if err != test.expectedErr || gauge != test.expectedValue {
				t.Errorf("\nValue: %f Expected: %f \nError: %s Expected: %s ", test.mValue, test.expectedValue, err, test.expectedErr)
			}
		})
	}
}

func TestServerCounter(t *testing.T) {
	cfg := config.Testcfg
	server := NewServerService(cfg)
	tests := []struct { // добавляем слайс тестов
		name          string
		mName         string
		mValue        int64
		expectedErr   error
		expectedValue int64
	}{
		{
			name:          "Gauge",
			mName:         "SomeGauge",
			mValue:        123,
			expectedValue: 123,
			expectedErr:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := server.UpdateCounter(test.mName, test.mValue)
			if err != test.expectedErr {
				t.Errorf("%s, Error: %s", test.name, err)
			}
			counter, err := server.ViewCounterByName(test.mName)
			if err != test.expectedErr || counter != test.expectedValue {
				t.Errorf("\nValue: %d Expected: %d \nError: %s Expected: %s ", test.mValue, test.expectedValue, err, test.expectedErr)
			}
		})
	}
}
