package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) UpdateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		name := vars["name"]
		value := vars["value"]
		if name == "" || value == "" {
			resp.WriteHeader(http.StatusBadRequest)
			http.Error(resp, "Metric value or name is doesn't exist", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(resp, req)
	})
}

func (s *Server) UpdateGauge(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	mValue, err := strconv.ParseFloat(vars["value"], 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.Service.UpdateGauge(vars["name"], mValue)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	} else {
		res.WriteHeader(http.StatusOK)
		return
	}
}

func (s *Server) UpdateCounter(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	mValue, err := strconv.ParseInt(vars["value"], 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

	err = s.Service.UpdateCounter(vars["name"], mValue)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func (s *Server) ValueMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		name := vars["name"]
		if name == "" {
			resp.WriteHeader(http.StatusBadRequest)
			http.Error(resp, "Missing metric value", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(resp, req)
	})
}

func (s *Server) ValueGauge(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metric, err := s.Service.ViewGaugeByName(vars["name"])
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(fmt.Sprint(metric)))
	} else {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(fmt.Sprint(metric)))
	}
}

func (s *Server) ValueCounter(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metric, err := s.Service.ViewCounterByName(vars["name"])
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(fmt.Sprint(metric)))
	} else {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(fmt.Sprintln(metric)))
	}
}

func (s *Server) MainPage(res http.ResponseWriter, req *http.Request) {
	list := s.Service.MetricList()
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(list))
}

func (s *Server) BadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
}
