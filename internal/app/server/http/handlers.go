package http

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggedResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggedResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggedResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

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

func (s *ServerApp) Logger(next http.Handler) http.Handler {

	loggedFunc := func(resp http.ResponseWriter, req *http.Request) {
		// start := time.Now()
		// method := req.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		loggedResp := loggedResponseWriter{
			ResponseWriter: resp, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		next.ServeHTTP(&loggedResp, req)
		// duration := time.Since(start)

		// reqLog := log.WithFields(log.Fields{
		// 	"URI":      req.RequestURI,
		// 	"method":   method,
		// 	"duration": duration,
		// })
		// reqLog.Info("request")

		respLog := log.WithFields(log.Fields{
			"status": responseData.status,
			"size":   responseData.size,
		})
		respLog.Info("response")

	}
	return http.HandlerFunc(loggedFunc)
}
