package service

import "errors"

var (
	serviceServiceError        = errors.New("service error")
	serviceStorageError        = errors.New("storage error")
	serviceConvertationError   = errors.New("string covertation error")
	serviceMetricNameIncorrect = errors.New("metric name is incorrect")
)
