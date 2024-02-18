package server

import (
	"fmt"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/models"
	"github.com/imirjar/metrx/internal/storage/database"
	"github.com/imirjar/metrx/internal/storage/mock"
)

func NewServerService(cfg config.ServerConfig) *ServerService {
	backupService := ServerService{
		cfg: cfg,
	}

	if err := backupService.PingDB(); err != nil {
		mock := mock.NewMockStorage(cfg)
		backupService.MemStorager = mock
		backupService.Backuper = mock
		fmt.Println("MEMSTORAGE")

	} else {
		db := database.NewDB(cfg)
		backupService.MemStorager = db
		backupService.Backuper = db
		fmt.Println("DB STORAGE")
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
	AddGauge(metric models.Metrics) error
	AddCounter(name string, delta int64) error
	ReadOne(metric models.Metrics) (models.Metrics, bool)
	ReadAllGauges() (map[string]float64, error)
	ReadAllCounters() (map[string]int64, error)
}

type DBManager interface {
	Import(path string) error
	Export(path string) error
}
