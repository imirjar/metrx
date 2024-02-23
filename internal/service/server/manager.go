package server

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"time"
// )

// func (s *ServerService) PeriodicBackup(interval time.Duration) {
// 	for {
// 		fmt.Println("PeriodicBackup")
// 		time.Sleep(interval)
// 		s.Backuper.Export(s.cfg.FilePath)
// 	}
// }

func (s *ServerService) Backup() error {
	if err := s.Backuper.Export(s.cfg.FilePath); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *ServerService) Restore() error {
	if err := s.Backuper.Import(s.cfg.FilePath); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
