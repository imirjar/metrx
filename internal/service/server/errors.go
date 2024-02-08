package server

import "errors"

var (
	errServiceError        = errors.New("service error")
	errStorageError        = errors.New("storage error")
	errConvertationError   = errors.New("string covertation error")
	errMetricNameIncorrect = errors.New("metric name is incorrect")
	errDBConnError         = errors.New("incorrect DB params")
)
