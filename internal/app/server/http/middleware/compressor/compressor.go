package compressor

import (
	"net/http"
	"strings"
)

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные

func Compressor(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		ow := resp
		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(resp)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(req.Body)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			req.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, req)
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
