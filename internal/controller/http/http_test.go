package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
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

type MockService struct {
}

func (fs *MockService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	return nil
}
func (fs *MockService) UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	return metric, nil
}
func (fs *MockService) ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	switch metric.MType {
	case "gauge":
		var val float64 = 100
		metric.Value = &val
	case "counter":
		var del int64 = 100
		metric.Delta = &del
	default:
		return metric, errMetricTypeIncorrect
	}

	return metric, nil
}
func (fs *MockService) MetricPage(ctx context.Context) (string, error) {
	return "", nil
}
