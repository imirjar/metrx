package service

import (
	"context"
)

func NewServerService() *ServerService {
	return &ServerService{}
}

type ServerService struct {
	MemStorager Storager
}

type Storager interface {
	AddGauges(ctx context.Context, gauges map[string]float64) error
	AddCounters(ctx context.Context, counters map[string]int64) error
	AddGauge(ctx context.Context, name string, value float64) (float64, error)
	AddCounter(ctx context.Context, name string, delta int64) (int64, error)
	ReadGauge(ctx context.Context, name string) (float64, bool)
	ReadCounter(ctx context.Context, name string) (int64, bool)
	ReadAllGauges(ctx context.Context) (map[string]float64, error)
	ReadAllCounters(ctx context.Context) (map[string]int64, error)
}
