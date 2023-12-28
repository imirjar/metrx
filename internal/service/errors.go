package service

import "errors"

var (
	serviceError        = errors.New("Service error")
	storageError        = errors.New("Storage error")
	convertationError   = errors.New("String covertation error")
	metricNameIncorrect = errors.New("Metric name is incorrect")
)
