package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *ServerApp) UpdateGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	mValue, err := strconv.ParseFloat(vars["value"], 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.Service.UpdateGauge(vars["name"], mValue)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}
}

func (s *ServerApp) UpdateCounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	mValue, err := strconv.ParseInt(vars["value"], 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
	}

	err = s.Service.UpdateCounter(vars["name"], mValue)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
	} else {
		resp.WriteHeader(http.StatusOK)
	}
}

func (s *ServerApp) UpdateBadParams(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
}

func (s *ServerApp) ValueGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metric, err := s.Service.ViewGaugeByName(vars["name"])
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(fmt.Sprint(metric)))
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprint(metric)))
	}
}

func (s *ServerApp) ValueCounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metric, err := s.Service.ViewCounterByName(vars["name"])
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(fmt.Sprint(metric)))
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprintln(metric)))
	}
}

func (s *ServerApp) ValueBadParams(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
}

func (s *ServerApp) MainPage(resp http.ResponseWriter, req *http.Request) {
	list := s.Service.MetricList()
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(list))
}
