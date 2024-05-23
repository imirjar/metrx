package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/imirjar/metrx/internal/models"
)

func URLParamsToMetric(r *http.Request) (models.Metrics, error) {

	var metric models.Metrics
	metric.ID = chi.URLParam(r, "name")
	metric.MType = chi.URLParam(r, "type")
	val := chi.URLParam(r, "value")
	if val != "" {
		if err := metric.SetVal(val); err != nil {
			log.Println("URLParams ERROR", err)
			return metric, err
		}
	}
	return metric, nil

}
