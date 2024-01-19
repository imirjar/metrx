package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) Update(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sName, sValue := vars["name"], vars["value"]
	switch vars["mType"] {
	case "gauge":
		mValue, err := strconv.ParseFloat(sValue, 64)
		if err != nil {
			fmt.Errorf("convertation error %s", err)
		}
		err = s.Service.UpdateGauge(sName, mValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	case "counter":
		mValue, err := strconv.ParseInt(sValue, 10, 64)
		if err != nil {
			fmt.Errorf("convertation error %s", err)
		}
		err = s.Service.UpdateCounter(sName, mValue)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
	}

}

func (s *Server) View(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	switch vars["mType"] {
	case "gauge":
		metric, err := s.Service.ViewGaugeByName(vars["name"])
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(fmt.Sprint(metric)))
		} else {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprint(metric)))
		}

	case "counter":
		metric, err := s.Service.ViewCounterByName(vars["name"])
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte(fmt.Sprint(metric)))
		} else {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(fmt.Sprintln(metric)))
		}

	default:
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Server) MainPage(res http.ResponseWriter, req *http.Request) {
	list := s.Service.MetricList()
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(list))
}
