package database

import "errors"

/* AddCounters(ctx context.Context, counters map[string]int64) error */
var (
	errAddCountersCloseError = errors.New("STORAGE AddCounters CLOSE ERROR")
)

/* AddGauges(ctx context.Context, gauges map[string]float64) error */
var (
	errAddGaugesCloseError = errors.New("STORAGE AddGauges CLOSE ERROR")
)

/* AddGauge(ctx context.Context, name string, value float64) (float64, error) */
var (
	errAddGaugeExecError = errors.New("STORAGE AddGauge EXEC ERROR")
)

/* AddCounter(ctx context.Context, name string, delta int64) (int64, error) */
var (
	errAddCounterExecError = errors.New("STORAGE AddCounter EXEC ERROR")
	errAddCounterScanError = errors.New("STORAGE AddCounter SCAN ERROR")
)

/* ReadAllGauges(ctx context.Context) (map[string]float64, error) */
var (
	errReadAllCountersQueryError = errors.New("STORAGE ReadAllCounters QUERY ERROR")
	errReadAllCountersRowsError  = errors.New("STORAGE ReadAllCounters ROWS ERROR")
	errReadAllCountersScanError  = errors.New("STORAGE ReadAllCounters SCAN ERROR")
)

/* ReadAllCounters(ctx context.Context) (map[string]int64, error) */
var (
	errReadAllGaugesQueryError = errors.New("STORAGE ReadAllGauges QUERY ERROR")
	errReadAllGaugesRowsError  = errors.New("STORAGE ReadAllGauges ROWS ERROR")
	errReadAllGaugesScanError  = errors.New("STORAGE ReadAllGauges SCAN ERROR")
)
