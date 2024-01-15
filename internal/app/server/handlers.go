package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imirjar/metrx/internal/service"
)

func newHandler() *Handler {
	return &Handler{
		service: service.New(),
	}
}

func (s *server) MainPage(res http.ResponseWriter, req *http.Request) {
	list := s.handler.service.MetricList()
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(list))
}

func (s *server) Update(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metricName, metricValue := vars["name"], vars["value"]
	switch vars["mType"] {
	case "gauge":
		err := s.handler.service.UpdateMetric("gauge", metricName, metricValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	case "counter":
		err := s.handler.service.UpdateMetric("counter", metricName, metricValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
	}

}

func (s *server) View(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	switch vars["mType"] {
	case "gauge":
		metric, err := s.handler.service.ViewGaugeByName(vars["name"])
		if err != nil {
			http.Error(res, "Metric type must be gauge or counter!", http.StatusBadRequest)
			return
		}
		if metric != nil {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintln(metric.Value)))
		} else {
			res.WriteHeader(http.StatusNotFound)
		}

	case "counter":
		metric, err := s.handler.service.ViewCounterByName(vars["name"])
		if err != nil {
			http.Error(res, "Metric type must be gauge or counter!", http.StatusBadRequest)
			return
		}
		if metric != nil {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintln(metric.Value)))
		} else {
			res.WriteHeader(http.StatusNotFound)
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}

func newRouter() *mux.Router {
	return mux.NewRouter()
}

type Handler struct {
	service *service.Service
}
