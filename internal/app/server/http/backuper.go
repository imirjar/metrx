package http

import (
	"net/http"
)

func (h *HttpApp) Backuper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		next.ServeHTTP(resp, req)
		h.Service.Storage.Export(h.DumpPath)
	})
}
