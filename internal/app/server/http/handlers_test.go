package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/app"
)

func TestServerApp_GaugeHandlers(t *testing.T) {
	serverApp := app.NewServerApp()
	type want struct {
		updateStatus int
		valueStatus  int
		value        string
	}
	type metricParams struct {
		value string
		name  string
	}
	tests := []struct {
		name   string
		params metricParams
		want   want
	}{
		{
			name: "#1 OK",
			params: metricParams{
				name:  "someGauge",
				value: "100",
			},
			want: want{
				updateStatus: 200,
				valueStatus:  200,
				value:        "100",
			},
		},
		{
			name: "#2 bad value",
			params: metricParams{
				name:  "none",
				value: "none",
			},
			want: want{
				updateStatus: 400,
				valueStatus:  404,
				value:        "0",
			},
		},
		{
			name: "#3 without value",
			params: metricParams{
				name:  "nil",
				value: "",
			},
			want: want{
				updateStatus: 404,
				valueStatus:  404,
				value:        "0",
			},
		},
		{
			name: "#4 resend 1 test",
			params: metricParams{
				name:  "someGauge",
				value: "100",
			},
			want: want{
				updateStatus: 200,
				valueStatus:  200,
				value:        "100",
			},
		},
	}

	for _, test := range tests {
		updateRecorder := httptest.NewRecorder()
		updatePath := fmt.Sprintf("/update/gauge/%s/%s", test.params.name, test.params.value)
		updateReq, err := http.NewRequest("POST", updatePath, nil)
		if err != nil {
			t.Fatal(err)
		}

		valueRecorder := httptest.NewRecorder()
		valuePath := fmt.Sprintf("/value/gauge/%s", test.params.name)
		valueReq, err := http.NewRequest("GET", valuePath, nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/update/gauge/{name}/{value}", serverApp.UpdateGauge)
		router.HandleFunc("/value/gauge/{name}", serverApp.ValueGauge)

		router.ServeHTTP(updateRecorder, updateReq)
		if updateRecorder.Code != test.want.updateStatus {
			t.Errorf("Error on %s update Gauge: status code %v want %v",
				test.name, updateRecorder.Code, test.want.updateStatus)
		}

		router.ServeHTTP(valueRecorder, valueReq)
		if valueRecorder.Body.String() != test.want.value || valueRecorder.Code != test.want.valueStatus {
			t.Errorf("Error on %s value Gauge: value >>%v<< want %v, status %d want %d",
				test.name, valueRecorder.Body.String(), test.want.value, valueRecorder.Code, test.want.valueStatus)
		}

	}
}

func TestServerApp_CounterHandlers(t *testing.T) {
	serverApp := app.NewServerApp()
	type want struct {
		updateStatus int
		valueStatus  int
		value        string
	}
	type metricParams struct {
		value string
		name  string
	}
	tests := []struct {
		name   string
		params metricParams
		want   want
	}{
		{
			name: "#1 OK",
			params: metricParams{
				name:  "someCounter",
				value: "100",
			},
			want: want{
				updateStatus: 200,
				valueStatus:  200,
				value:        "100",
			},
		},
		{
			name: "#2 bad value",
			params: metricParams{
				name:  "none",
				value: "none",
			},
			want: want{
				updateStatus: 400,
				valueStatus:  404, //err
				value:        "0",
			},
		},
		{
			name: "#3 without value",
			params: metricParams{
				name:  "nil",
				value: "",
			},
			want: want{
				updateStatus: 404,
				valueStatus:  404,
				value:        "0",
			},
		},
		{
			name: "#4 resend 1 test",
			params: metricParams{
				name:  "someCounter",
				value: "100",
			},
			want: want{
				updateStatus: 200,
				valueStatus:  200,
				value:        "200",
			},
		},
	}

	for _, test := range tests {
		updateRecorder := httptest.NewRecorder()
		updatePath := fmt.Sprintf("/update/counter/%s/%s", test.params.name, test.params.value)
		updateReq, err := http.NewRequest("POST", updatePath, nil)
		if err != nil {
			t.Fatal(err)
		}

		valueRecorder := httptest.NewRecorder()
		valuePath := fmt.Sprintf("/value/counter/%s", test.params.name)
		valueReq, err := http.NewRequest("GET", valuePath, nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/update/counter/{name}/{value}", serverApp.UpdateCounter)
		router.HandleFunc("/value/counter/{name}", serverApp.ValueCounter)

		router.ServeHTTP(updateRecorder, updateReq)
		if updateRecorder.Code != test.want.updateStatus {
			t.Errorf("Error on %s update Counter: status code %v want %v",
				test.name, updateRecorder.Code, test.want.updateStatus)
		}

		router.ServeHTTP(valueRecorder, valueReq)
		if valueRecorder.Body.String() != test.want.value {
			t.Errorf("Error on %s value Counter: value >>%v<< want %v",
				test.name, valueRecorder.Body.String(), test.want.value)
		}
		if valueRecorder.Code != test.want.valueStatus {
			t.Errorf("Error on %s value Counter: status %d want %d",
				test.name, valueRecorder.Code, test.want.valueStatus)
		}

	}
}
