package database

import (
	"database/sql"

	"github.com/imirjar/metrx/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func NewDB(cfg config.ServerConfig) *DB {
	conn, err := sql.Open("pgx", cfg.DBConn)
	if err != nil {
		panic(err)
	}
	storage := DB{
		db: conn,
	}

	// storage.configure(cfg)

	return &storage
}

// func (m *DB) up() error {
// 	_, err := m.db.Exec(
// 		`CREATE TABLE metrics (
// 			"id" VARCHAR(50) NOT NULL
// 			"type" VARCHAR(250) NOT NULL,
// 			"delta" FLOAT,
// 			"value" INTEGER
// 		) `)
// 	return err
// }

// func (m *DB) down() error {
// 	_, err := m.db.Exec("DROP TABLE metrics")
// 	return err
// }
