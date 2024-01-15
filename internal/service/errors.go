package service

import "errors"

var (
	serviceError        = errors.New("service error")
	storageError        = errors.New("storage error")
	convertationError   = errors.New("string covertation error")
	metricNameIncorrect = errors.New("metric name is incorrect")
)
