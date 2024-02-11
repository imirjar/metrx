package backup

import (
	"context"
	"log"
	"time"

	"github.com/imirjar/metrx/config"
	"github.com/imirjar/metrx/internal/storage/mock"
)

func NewBackupService(cfg config.ServerConfig) *BackupService {
	store := mock.NewMockStorage(&cfg.StorageConfig)
	backupService := BackupService{
		Backuper: store,
		cfg:      &cfg.ServiceConfig,
	}
	// run dump auto-exporter
	if backupService.cfg.Interval > 0 {
		go backupService.PeriodicBackup(backupService.cfg.Interval)
	}
	return &backupService
}

type BackupService struct {
	Backuper Backuper
	cfg      *config.ServiceConfig
}

type Backuper interface {
	Import() error
	Export() error
	Ping(ctx context.Context) (bool, error)
}

func (s *BackupService) CheckDBConn(ctx context.Context) (bool, error) {
	return s.Backuper.Ping(ctx)
	// if s.cfg.DBConn == "" {
	// 	return false, errDBConnError
	// }

	// db, err := sql.Open("pgx", s.Cfg.DBConn)
	// if err != nil {
	// 	// log.Fatalf(err.Error())
	// 	return false, err
	// }

	// err = db.PingContext(ctx)
	// if err != nil {
	// 	// log.Fatalf(err.Error())
	// 	return false, err
	// }
	// return true, nil
}

// Backuper
// loop backup
func (s *BackupService) PeriodicBackup(interval time.Duration) {
	for {
		time.Sleep(interval)
		s.Backuper.Export()
	}
}

func (s *BackupService) Backup() error {
	if err := s.Backuper.Export(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *BackupService) Restore() error {
	if err := s.Backuper.Import(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
