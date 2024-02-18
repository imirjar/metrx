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
	Create(metric models.Metrics) error
	ReadOne(metric models.Metrics) (models.Metrics, bool)
	ReadAll(mType string) ([]models.Metrics, error)
	Update(metric models.Metrics) error
	Delete(metric models.Metrics) error
	// AddGauge(mName string, mValue float64) (float64, error)
	// AddCounter(mName string, mValue int64) (int64, error)
	// ReadAllGauge() map[string]float64
	// ReadAllCounter() map[string]int64
	// ReadGauge(mName string) (float64, bool)
	// ReadCounter(mName string) (int64, bool)
}

type DBManager interface {
	Import(path string) error
	Export(path string) error
}
