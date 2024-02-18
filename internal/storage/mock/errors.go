package mock

import "errors"

var (
	errDBConnError          = errors.New("incorrect DB params")
	errMetricStructureError = errors.New("incorrect DB params")
)
