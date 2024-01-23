package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *ServerApp) UpdateGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	vv, ok := vars["value"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mValue, err := strconv.ParseFloat(vv, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.Service.UpdateGauge(vn, mValue)
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

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	vv, ok := vars["value"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mValue, err := strconv.ParseInt(vv, 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.Service.UpdateCounter(vn, mValue)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}
}

func (s *ServerApp) ValueGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := s.Service.ViewGaugeByName(vn)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(fmt.Sprint(metric)))
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprint(metric)))
		return
	}
}

func (s *ServerApp) ValueCounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := s.Service.ViewCounterByName(vn)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(fmt.Sprint(metric)))
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprint(metric)))
	}
}

func (s *ServerApp) MainPage(resp http.ResponseWriter, req *http.Request) {
	list := s.Service.MetricList()
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(list))
}

func (s *ServerApp) BadParams(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
}
