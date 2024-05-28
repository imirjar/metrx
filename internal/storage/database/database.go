package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	db *pgx.Conn
}

func NewDB(DBConn string) *DB {
	conn, err := pgx.Connect(context.Background(), DBConn)
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
