package server

import (
	"net/http"
)

func Run() error {
	server := NewServer()

	return http.ListenAndServe(server.Config.url, server.Router)
}
