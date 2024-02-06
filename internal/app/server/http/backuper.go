package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (h *HTTPApp) Backuper(next http.Handler) http.Handler {
	err := h.Service.Backup()
	if err != nil {
		log.Error(err)
		return http.HandlerFunc(h.Error)
	}
	return next
	// next.ServeHTTP(resp, req)
}
