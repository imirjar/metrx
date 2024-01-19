package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imirjar/metrx/internal/app/server"
)

func Test_server_Update_view(t *testing.T) {
	s := server.NewServer()
	w := httptest.NewRecorder()

	s.Router.HandleFunc("/update/{mType}/{name}/{value}", s.Update)
	s.Router.ServeHTTP(w, httptest.NewRequest("POST", "/update/gauge/someGauge/1", nil))

	if w.Code != http.StatusOK {
		t.Error("Did not get expected HTTP status code, got", w.Code)
	}
}
