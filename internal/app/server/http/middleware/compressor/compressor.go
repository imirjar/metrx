package compressor

import (
	"net/http"
	"strings"
)

func Compressor(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		ow := resp
		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(resp)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}

			req.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, req)

		// for question

		// encodingHeaders := req.Header.Values("Content-Encoding")
		// if !slices.Contains(encodingHeaders, "gzip") {
		// 	next.ServeHTTP(resp, req)
		// 	return
		// }

		// cmpReader, err := gzip.NewReader(req.Body)
		// if err != nil {
		// 	http.Error(resp, err.Error(), http.StatusNotFound)
		// 	return
		// }
		// defer cmpReader.Close()
		// req.Body = cmpReader

		// dcmpWriter := gzip.NewWriter(resp)
		// defer dcmpWriter.Close()

		// next.ServeHTTP(gzipWriter{ResponseWriter: resp, Writer: dcmpWriter}, req)

	})
}
