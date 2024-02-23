package database

import (
	"context"
	"log"

	"github.com/imirjar/metrx/config"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	db *pgx.Conn
}

func NewDB(cfg config.ServerConfig) *DB {
	conn, err := pgx.Connect(context.Background(), cfg.DBConn)

	if err != nil {
		log.Fatalln(err)
	}
	storage := DB{
		db: conn,
	}

	err = storage.Up()
	if err != nil {
		log.Fatalln(err)
	}

	// storage.configure(cfg)

	return &storage
}

func (m *DB) Up() error {
	_, err := m.db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS metrics (
			id varchar NOT NULL,
			"type" varchar NOT NULL,
			value float8 NULL,
			CONSTRAINT metrics_pk PRIMARY KEY (id)
		);`)
	return err
}

// func (m *DB) Down() error {
// 	_, err := m.db.Exec("DROP TABLE metrics")
// 	return err
// }
