package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HTTPGateway) Ping(resp http.ResponseWriter, req *http.Request) {
	err := h.Service.PingDB()
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(err.Error()))
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
}

func (h *HTTPGateway) Update(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok || vn == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}

	vt, ok := vars["type"]
	if !ok || !(vt == "gauge" || vt == "counter") {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
		return
	}

	vv, ok := vars["value"]
	if !ok || vv == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusBadRequest)
		return
	}

	err := h.Service.Update(vn, vt, vv) //here
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}
}

func (h *HTTPGateway) View(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok || vn == "" {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
		return
	}

	vt, ok := vars["type"]
	if !ok || !(vt == "gauge" || vt == "counter") {
		http.Error(resp, errMetricNameIncorrect.Error(), http.StatusNotFound)
		return
	}

	value, err := h.Service.View(vn, vt) //here
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(value))
		return
	}
}

func (h *HTTPGateway) MainPage(resp http.ResponseWriter, req *http.Request) {
	page := h.Service.MetricPage()
	resp.Header().Set("content-type", "text/html")
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, page)
}

func (h *HTTPGateway) UpdateJSON(resp http.ResponseWriter, req *http.Request) {
	// var metric models.Metrics
	var buf bytes.Buffer //byte buffer

	//request body -> buf
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	bMetric, err := h.Service.ByteUpdate(buf.Bytes())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(bMetric)
}

func (h *HTTPGateway) ValueJSON(resp http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer //byte buffer

	//request body -> buf
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	bMetric, err := h.Service.ByteRead(buf.Bytes())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Set("content-type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(bMetric)
}
