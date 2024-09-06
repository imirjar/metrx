package http

import "errors"

var (
	errMetricNameIncorrect = errors.New("metric name is incorrect")
	errParamsIncorrect     = errors.New("metric params is incorrect")

	errLoadPrivateKey  = errors.New("can't load private key")
	errBlockPrivateKey = errors.New("can't block private key")
)
