package server

import (
	"context"
	"log"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/storage/database"
	"github.com/imirjar/metrx/internal/storage/mock"
	"github.com/imirjar/metrx/pkg/ping"
)

func NewServerService(cfg config.ServerConfig) *ServerService {
	backupService := ServerService{
		cfg: cfg,
	}

	db, err := ping.NewDBPool(context.Background(), cfg.DBConn)
	if err != nil {
		log.Print(err)
	}
	if err = db.Ping(context.Background()); err != nil {
		mock := mock.NewMockStorage(cfg)
		backupService.MemStorager = mock

	} else {
		db := database.NewDB(cfg)
		backupService.MemStorager = db
	}

	return &backupService
}

type ServerService struct {
	MemStorager Storager
	cfg         config.ServerConfig
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
