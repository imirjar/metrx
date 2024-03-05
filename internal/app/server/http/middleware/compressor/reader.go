package compressor

import (
	"compress/gzip"
	"io"
)

type compressReader struct {
	req   io.ReadCloser
	gzReq *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		req:   r,
		gzReq: zr,
	}, nil
}
