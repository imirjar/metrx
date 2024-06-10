package models

import "errors"

var (
	errMetricTypeError    = errors.New("incorrect metric type")
	errTypeAssertionError = errors.New("failed to convert value to type")
)
