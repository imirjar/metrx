package compressor

import (
	"compress/gzip"
	"io"
)

type compressReader struct {
	req   io.ReadCloser
	gzReq *gzip.Reader
}

func newCompressReader(req io.ReadCloser) (*compressReader, error) {
	gzReq, err := gzip.NewReader(req)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		req:   req,
		gzReq: gzReq,
	}, nil
}
