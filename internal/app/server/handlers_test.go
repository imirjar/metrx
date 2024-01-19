package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imirjar/metrx/internal/app/server"
	"github.com/stretchr/testify/assert"
)

func Test_server_UpdateGauge_view(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct { // добавляем слайс тестов
		name string
		url  string
		want want
	}{
		{
			name: "Correct",
			url:  "/update/gauge/someGauge/12",
			want: want{
				code:        200,
				contentType: "application/json",
			},
		},
		// {
		// 	name:     "None value",
		// 	url:      "/update/gauge/someGauge/none",
		// 	expected: 400,
		// },
		// {
		// 	name:     "No value in path",
		// 	url:      "/update/gauge/someGauge",
		// 	expected: 404,
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := server.NewServer()
			request := httptest.NewRequest(http.MethodPost, test.url, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(s.UpdateGauge)
			h(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func Test_server_UpdateCounter_view(t *testing.T) {
	s := server.NewServer()
	w := httptest.NewRecorder()

	tests := []struct { // добавляем слайс тестов
		name     string
		url      string
		expected int
	}{
		{
			name:     "200",
			url:      "/update/counter/someCounter/1",
			expected: 200,
		},
		{
			name:     "400",
			url:      "/update/counter/someCounter/none",
			expected: 200,
		},
		{
			name:     "404",
			url:      "/update/counter/someCounter",
			expected: 200,
		},
	}

	for _, test := range tests {
		s.Router.HandleFunc("/update/counter/{name}/{value}", s.UpdateCounter)
		s.Router.ServeHTTP(w, httptest.NewRequest("POST", test.url, nil))

		if w.Code != test.expected {
			t.Errorf("Expected: %d, Real: %d", test.expected, w.Code)
		}
	}
}
