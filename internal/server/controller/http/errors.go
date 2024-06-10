package http

import "errors"

var (
	errMetricNameIncorrect  = errors.New("metric name is incorrect")
	errParamsIncorrect      = errors.New("metric params is incorrect")
	errMetricTypeIncorrect  = errors.New("metric type is incorrect")
	errMetricValueIncorrect = errors.New("metric value is incorrect")
)
