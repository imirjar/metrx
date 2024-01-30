package pkg

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = zw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return buf.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	// переменная r будет читать входящие данные и распаковывать их

	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	var b bytes.Buffer
	// в переменную b записываются распакованные данные
	_, err = b.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
