package mock

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
)

func (m *MemStorage) Export() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	file, err := os.OpenFile(m.cfg.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	file.Write(data)

	return nil
}

func (m *MemStorage) Import() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	file, err := os.ReadFile(m.cfg.FilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, m); err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) Ping(ctx context.Context) (bool, error) {
	if m.cfg.DBConn == "" {
		return false, errDBConnError
	}

	db, err := sql.Open("pgx", m.cfg.DBConn)
	if err != nil {
		// log.Fatalf(err.Error())
		return false, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		// log.Fatalf(err.Error())
		return false, err
	}
	return true, nil
}
