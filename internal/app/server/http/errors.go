package http

import "errors"

var (
	errMetricNameIncorrect  = errors.New("metric name is incorrect")
	errMetricTypeUnexpected = errors.New("metric type is incorrect")
)
