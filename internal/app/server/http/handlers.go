package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/imirjar/metrx/internal/entity"

	"github.com/gorilla/mux"
)

func (h *HTTPApp) Ping(resp http.ResponseWriter, req *http.Request) {

	_, err := h.Service.CheckDBConn(req.Context())
	if err != nil {
		http.Error(resp, "Failed DB connection", http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))

}

func (h *HTTPApp) UpdateGauge(resp http.ResponseWriter, req *http.Request) {
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
	err = h.Service.UpdateGauge(vn, mValue)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}
}

func (h *HTTPApp) UpdateCounter(resp http.ResponseWriter, req *http.Request) {
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
	err = h.Service.UpdateCounter(vn, mValue)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	} else {
		resp.WriteHeader(http.StatusOK)
		return
	}
}

func (h *HTTPApp) ValueGauge(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.Service.ViewGaugeByName(vn)
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

func (h *HTTPApp) ValueCounter(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	vn, ok := vars["name"]
	if !ok {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.Service.ViewCounterByName(vn)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(fmt.Sprint(metric)))
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(fmt.Sprint(metric)))
	}
}

func (h *HTTPApp) MainPage(resp http.ResponseWriter, req *http.Request) {
	page := h.Service.MetricPage()
	resp.Header().Set("content-type", "text/html")
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, page)
}

func (h *HTTPApp) UpdateJSON(resp http.ResponseWriter, req *http.Request) {

	var metric entity.Metrics
	var buf bytes.Buffer //byte buffer

	//request body -> buf
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	//put body JSON to metric model
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	//counter||gauge||StatusNotFound
	switch metric.MType {
	case "gauge":
		//update value in storage
		h.Service.UpdateGauge(metric.ID, *metric.Value) //надо возвращать обновленное значение!

		r, err := json.Marshal(metric)
		if err != nil {
			//400 if marshal error
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		//200 with Content-type:application/json
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)

	case "counter":
		//update value in storage
		h.Service.UpdateCounter(metric.ID, *metric.Delta)

		r, err := json.Marshal(metric)
		if err != nil {
			//400 if marshal error
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		//200 with Content-type:application/json
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)

	default:
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write(nil)
	}
}

func (h *HTTPApp) ValueJSON(resp http.ResponseWriter, req *http.Request) {

	var metric entity.Metrics
	var buf bytes.Buffer //byte buffer

	//request body -> buf
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	//put body JSON to metric model
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	//counter||gauge||StatusNotFound
	switch metric.MType {
	case "gauge":
		//get value from storage
		v, err := h.Service.ViewGaugeByName(metric.ID)
		if err != nil {
			//404 if metric doesn't found in storage
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusNotFound)
			resp.Write(nil)
		}
		metric.Value = &v

		r, err := json.Marshal(metric)
		if err != nil {
			//400 if marshal error
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write(nil)
		}
		//200 with Content-type:application/json
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)
	case "counter":
		v, err := h.Service.ViewCounterByName(metric.ID)
		if err != nil {
			//404 if metric doesn't found in storage
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusNotFound)
			resp.Write(nil)
		}
		metric.Delta = &v

		r, err := json.Marshal(metric)
		if err != nil {
			//400 if marshal error
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write(nil)
		}
		//200 with Content-type:application/json
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)
	default:
		//404 for all values
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write(nil)
	}
}

func (h *HTTPApp) BadParams(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
}

func (h *HTTPApp) Error(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusInternalServerError)
}

// func (s *ServerApp) error(w http.ResponseWriter, r *http.Request, code int, err error) {
// 	s.respond(w, r, code, map[string]string{"error": err.Error()})
// }

// func (s *ServerApp) respond(resp http.ResponseWriter, req *http.Request, code int, data interface{}) {
// 	resp.Header().Set("content-type", "application/json")
// 	resp.WriteHeader(code)
// 	if data != nil {
// 		json.NewEncoder(resp).Encode(data)
// 	}

// }
