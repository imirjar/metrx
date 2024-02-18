package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/config"
	server "github.com/imirjar/metrx/internal/app/server/http"
)

func TestServerApp_GaugeHandlers(t *testing.T) {
	cfg := config.Testcfg
	serverApp := server.NewGateway(cfg)

	tests := []struct {
		tName          string
		mName          string
		mType          string
		mValue         string
		statusExpected int
	}{
		{
			tName:          "#1 OK gauge",
			mType:          "gauge",
			mName:          "someGauge",
			mValue:         "100",
			statusExpected: 200,
		},
		{
			tName:          "#1 OK counter",
			mType:          "counter",
			mName:          "testSetGet144",
			mValue:         "835",
			statusExpected: 200,
		},
	}

	for _, test := range tests {
		updateRecorder := httptest.NewRecorder()
		updatePath := fmt.Sprintf("/update/%s/%s/%s", test.mType, test.mName, test.mValue)
		updateReq, err := http.NewRequest("POST", updatePath, nil)
		if err != nil {
			t.Fatal(err)
		}

		valueRecorder := httptest.NewRecorder()
		valuePath := fmt.Sprintf("/value/%s/%s", test.mType, test.mName)
		valueReq, err := http.NewRequest("GET", valuePath, nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/update/gauge/{name}/{value}", serverApp.UpdateGauge)
		router.HandleFunc("/update/counter/{name}/{value}", serverApp.UpdateCounter)
		router.HandleFunc("/value/gauge/{name}", serverApp.ValueGauge)
		router.HandleFunc("/value/counter/{name}", serverApp.ValueCounter)

		router.ServeHTTP(updateRecorder, updateReq)
		if updateRecorder.Code != test.statusExpected {
			t.Errorf("Error on %s update Gauge: status code %v want %v",
				test.tName, updateRecorder.Code, test.statusExpected)
		}

		router.ServeHTTP(valueRecorder, valueReq)
		// valueRecorder.Body.String() != test.mValue
		if valueRecorder.Code != test.statusExpected {
			t.Errorf("Error on %s value Gauge: status %d want %d",
				test.tName, valueRecorder.Code, test.statusExpected)
		}

	}
}
