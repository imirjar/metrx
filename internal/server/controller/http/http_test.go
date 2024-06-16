package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/imirjar/metrx/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathUpdateHandler(t *testing.T) {
	type metric struct {
		name  string
		mtype string
		value string
	}

	type want struct {
		contentType string
		statusCode  int
		expected    string
	}

	tests := []struct {
		name   string
		want   want
		metric metric
	}{
		{
			name: "ok",
			want: want{
				contentType: "",
				statusCode:  http.StatusOK,
				expected:    "100",
			},
			metric: metric{
				mtype: "gauge",
				name:  "testGauge",
				value: "100",
			},
		},
		{
			name: "400",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				expected:    "error metric isn't float64\n",
			},
			metric: metric{
				mtype: "gauge",
				name:  "testGauge",
				value: "",
			},
		},
	}

	for _, tt := range tests {

		// UpdatePathHandler
		t.Run(tt.name, func(t *testing.T) {

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.metric.mtype)
			rctx.URLParams.Add("name", tt.metric.name)
			rctx.URLParams.Add("value", tt.metric.value)

			r := httptest.NewRequest(http.MethodPost, "/update/{type}/{name}/{value}", nil)
			w := httptest.NewRecorder()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(ts.UpdatePathHandler())
			h(w, r)

			res := w.Result()

			if ok := assert.Equal(t, tt.want.statusCode, w.Code); !ok {
				log.Print(r.URL.Host)
				t.Error(w.Body)
			}
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.expected, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})

	}
}

func TestPathValueHandler(t *testing.T) {
	type metric struct {
		name  string
		mtype string
	}

	type want struct {
		contentType string
		statusCode  int
		expected    string
	}

	tests := []struct {
		name   string
		want   want
		metric metric
	}{
		{
			name: "ok",
			want: want{
				contentType: "",
				statusCode:  http.StatusOK,
				expected:    "100",
			},
			metric: metric{
				mtype: "gauge",
				name:  "testGauge",
			},
		},
		{
			name: "404 wrong type",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				expected:    "error metric isn't float64\n",
			},
			metric: metric{
				mtype: "gauges",
				name:  "testGauge",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.metric.mtype)
			rctx.URLParams.Add("name", tt.metric.name)

			r := httptest.NewRequest(http.MethodPost, "/value/{type}/{name}", nil)
			w := httptest.NewRecorder()
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h := http.HandlerFunc(ts.ValuePathHandler())
			h(w, r)

			res := w.Result()

			if ok := assert.Equal(t, tt.want.statusCode, w.Code); !ok {
				log.Print(r.URL.Host)
				t.Error(w.Body)
			}
			defer res.Body.Close()

		})
	}
}

var ts HTTPGateway = HTTPGateway{
	Service: &MockService{},
}

func TestMainPageHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		someFloat64 float64 = 10
		someInt64   int64   = 100
	)

	type testCase struct {
		name           string
		mockService    func(service *MockService)
		expectedStatus int
		expectedBody   string
		expectedError  error
	}

	tests := []testCase{
		{
			name: "Success case",
			mockService: func(service *MockService) {
				expectedMetrics := map[string][]models.Metrics{
					"gauges": {
						{ID: "metric1", MType: "gauge", Value: &someFloat64},
						{ID: "metric2", MType: "gauge", Value: &someFloat64},
					},
					"counters": {
						{ID: "metric3", MType: "counter", Delta: &someInt64},
					},
				}
				service.EXPECT().ViewMetrics(gomock.Any()).Return(expectedMetrics, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "value1value2value3", // Проверьте, что HTML содержит ожидаемые значения
			expectedError:  nil,
		},
		// {
		// 	name: "Error case",
		// 	mockService: func(service *MockService) {
		// 		expectedError := errors.New("error fetching metrics")
		// 		service.EXPECT().ViewMetrics(gomock.Any()).Return(nil, expectedError)
		// 	},
		// 	expectedStatus: http.StatusInternalServerError,
		// 	expectedBody:   "error fetching metrics",
		// 	expectedError:  nil,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockService(ctrl)
			h := &HTTPGateway{
				Service: mockService,
				Secret:  "",
			}

			if tt.mockService != nil {
				tt.mockService(mockService)
			}

			req := httptest.NewRequest("GET", "/", nil)
			log.Print(req.Body)
			w := httptest.NewRecorder()

			handler := h.MainPage()
			handler(w, req)

			// log.Print(tt.expectedStatus)
			log.Print(w.Body)
			require.Equal(t, tt.expectedStatus, w.Code, "Status code mismatch")

			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody, "Response body mismatch")
			}

			if tt.expectedError != nil {
				assert.Contains(t, w.Body.String(), tt.expectedError.Error(), "Error message mismatch")
			}
		})
	}
}

// func (fs *MockService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
// 	return nil
// }
// func (fs *MockService) UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
// 	return metric, nil
// }
// func (fs *MockService) ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
// 	switch metric.MType {
// 	case "gauge":
// 		var val float64 = 100
// 		metric.Value = &val
// 	case "counter":
// 		var del int64 = 100
// 		metric.Delta = &del
// 	default:
// 		return metric, errMetricTypeIncorrect
// 	}

// 	return metric, nil
// }
// func (fs *MockService) ViewMetrics(ctx context.Context) (map[string][]models.Metrics, error) {
// 	return make(map[string][]models.Metrics), nil
// }
