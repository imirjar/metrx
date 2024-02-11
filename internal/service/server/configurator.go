package server

import (
	"context"
	"log"
	"time"
)

type Backuper interface {
	Import() error
	Export() error
	Ping(ctx context.Context) (bool, error)
}

func (s *ServerService) CheckDBConn(ctx context.Context) (bool, error) {
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
func (s *ServerService) PeriodicBackup(interval time.Duration) {
	for {
		time.Sleep(interval)
		s.Backuper.Export()
	}
}

func (s *ServerService) Backup() error {
	if err := s.Backuper.Export(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *ServerService) Restore() error {
	if err := s.Backuper.Import(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
