package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/imirjar/metrx/internal/entity"
	log "github.com/sirupsen/logrus"

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

func (s *ServerApp) UpdateJSON(resp http.ResponseWriter, req *http.Request) {
	var metric entity.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	switch metric.MType {
	case "gauge":
		s.Service.UpdateGauge(metric.ID, *metric.Value) //надо возвращать обновленное значение!
		r, err := json.Marshal(metric)
		if err != nil {
			//может быть проблемма связанная с тем что значение не найдено а не сбой
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)

	case "counter":
		s.Service.UpdateCounter(metric.ID, *metric.Delta)
		r, err := json.Marshal(metric)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)

	default:
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write(nil)
	}
}

func (s *ServerApp) ValueJSON(resp http.ResponseWriter, req *http.Request) {

	var metric entity.Metrics

	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println("###", metric, "###")

	switch metric.MType {
	case "gauge":
		v, err := s.Service.ViewGaugeByName(metric.ID)
		if err != nil {
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusNotFound)
			resp.Write(nil)
		}
		metric.Value = &v

		r, err := json.Marshal(metric)
		if err != nil {
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write(nil)
		}
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)
	case "counter":
		v, err := s.Service.ViewCounterByName(metric.ID)
		if err != nil {
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusNotFound)
			resp.Write(nil)
		}
		metric.Delta = &v
		r, err := json.Marshal(metric)
		if err != nil {
			resp.Header().Set("content-type", "application/json")
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write(nil)
		}
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write(r)
	default:
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write(nil)
	}

}

func (s *ServerApp) BadParams(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusBadRequest)
}

func (s *ServerApp) Logger(next http.Handler) http.Handler {

	loggedFunc := func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()
		method := req.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		loggedResp := loggedResponseWriter{
			ResponseWriter: resp, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		next.ServeHTTP(&loggedResp, req)
		duration := time.Since(start)

		reqLog := log.WithFields(log.Fields{
			"URI":      req.RequestURI,
			"method":   method,
			"duration": duration,
		})
		reqLog.Info("request")

		respLog := log.WithFields(log.Fields{
			"status": responseData.status,
			"size":   responseData.size,
		})
		respLog.Info("response")

	}
	return http.HandlerFunc(loggedFunc)
}

// func (s *ServerApp) error(w http.ResponseWriter, r *http.Request, code int, err error) {
// 	s.respond(w, r, code, map[string]string{"error": err.Error()})
// }

// func (s *ServerApp) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
// 	w.WriteHeader(code)
// 	if data != nil {
// 		json.NewEncoder(w).Encode(data)
// 	}
// }
