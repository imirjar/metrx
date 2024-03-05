package logger

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
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

func Logger(next http.Handler) http.Handler {

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
