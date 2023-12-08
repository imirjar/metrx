package storage

import (
	"testing"

	"github.com/imirjar/metrx/internal/models"
)

func TestGaugeCreate(t *testing.T) {

	store := New()

	gauge := models.Gauge{
		Name:  "Name",
		Value: 2,
	}

	store.GaugeCreate(&gauge)

	if store.GaugeStorage[0] != gauge {
		t.Errorf("sum expected to be 3; got %d", gauge)
	}
}
