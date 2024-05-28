package http

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/imirjar/metrx/internal/models"
)

var (
	fs                     = fakeService{}
	testServer HTTPGateway = HTTPGateway{
		Service: fs,
	}
)

func TestUpdatePathHandler(t *testing.T) {
	gw := HTTPGateway{}
	fg := httptest.NewServer(gw.UpdatePathHandler())
	defer fg.Close()
}

type fakeService struct {
}

func (fs fakeService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	return nil
}
func (fs fakeService) UpdateMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	return models.Metrics{}, nil
}
func (fs fakeService) ViewMetric(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	return models.Metrics{}, nil
}
func (fs fakeService) MetricPage(ctx context.Context) (string, error) {
	return "", nil
}
