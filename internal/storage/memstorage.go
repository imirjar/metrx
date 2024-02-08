package storage

import "database/sql"

type MemStorageDB struct {
	db *sql.DB
}

func NewMemStorageDB(db *sql.DB) *MemStorageDB {
	return &MemStorageDB{
		db: db,
	}
}
