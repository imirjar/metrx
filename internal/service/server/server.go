package server

import (
	"context"
	"log"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/storage/database"
	"github.com/imirjar/metrx/internal/storage/mock"
	"github.com/imirjar/metrx/pkg/ping"
)

func NewServerService(cfg config.ServerConfig) *ServerService {
	backupService := ServerService{
		cfg: cfg,
	}

	log.Println(cfg.DBConn)
	connIsValid := ping.PingPgx(context.Background(), cfg.DBConn)

	if connIsValid {
		db := database.NewDB(cfg)
		backupService.MemStorager = db
		backupService.Backuper = db
		// fmt.Println("DB STORAGE")

	} else {
		mock := mock.NewMockStorage(cfg)
		backupService.MemStorager = mock
		backupService.Backuper = mock
		// fmt.Println("MEMSTORAGE")
	}

	// store := mock.NewMockStorage(cfg)

	// run dump auto-exporter
	// if backupService.cfg.Interval > 0 {
	// 	go backupService.PeriodicBackup(backupService.cfg.Interval)
	// }

	return &backupService
}

type ServerService struct {
	MemStorager Storager
	Backuper    DBManager
	cfg         config.ServerConfig
}

type Storager interface {
	AddGauges(gauges map[string]float64) error
	AddCounters(counters map[string]int64) error
	AddGauge(name string, value float64) (float64, error)
	AddCounter(name string, delta int64) (int64, error)
	ReadOne(metric models.Metrics) (models.Metrics, bool)
	ReadAllGauges() (map[string]float64, error)
	ReadAllCounters() (map[string]int64, error)
}

type DBManager interface {
	Import(path string) error
	Export(path string) error
}
