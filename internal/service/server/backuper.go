package server

import (
	"log"
	"time"
)

type Backuper interface {
	Import() error
	Export() error
}

// Backuper
// loop backup
func (s *ServerService) PeriodicBackup(interval time.Duration) {

	for {
		time.Sleep(interval)
		s.Dump.Export()
	}
}

func (s *ServerService) Backup() error {
	if err := s.Dump.Export(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *ServerService) Restore() error {
	if err := s.Dump.Import(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
