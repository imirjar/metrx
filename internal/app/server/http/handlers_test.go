package http

import (
	"net/http"
	"testing"

	"github.com/imirjar/metrx/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestHTTPGateway_ValueJSONHandler(t *testing.T) {
	type fields struct {
		Service Service
		cfg     config.ServerConfig
	}
	type args struct {
		resp http.ResponseWriter
		req  *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPGateway{
				Service: tt.fields.Service,
				cfg:     tt.fields.cfg,
			}
			h.ValueJSONHandler(tt.args.resp, tt.args.req)
		})
	}
}
