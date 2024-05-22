package compressor

import (
	"compress/gzip"
	"io"
)

type CompressReader struct {
	Req   io.ReadCloser
	GzReq *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{
		Req:   r,
		GzReq: zr,
	}, nil
}
