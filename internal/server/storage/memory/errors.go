package memory

import "errors"

var (
	errDBConnError     = errors.New("incorrect DB params")
	errDBWriteError    = errors.New("unable to write in db")
	errMetricTypeError = errors.New("metric type is unexpected")
)
