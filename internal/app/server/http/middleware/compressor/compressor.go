package compressor

import (
	"net/http"
	"strings"
)

func Compressor(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		acceptEncoding := req.Header.Get("Accept-Encoding")
		contentEncoding := req.Header.Get("Content-Encoding")

		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if supportsGzip {
			cResp := newCompressWriter(resp)
			defer cResp.Close()
			resp = cResp
		}

		if sendsGzip {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer cr.Close()
			req.Body = cr
		}

		next.ServeHTTP(resp, req)
	})
}
