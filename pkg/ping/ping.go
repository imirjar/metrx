package ping

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDBPool(ctx context.Context, connString string) (*DB, error) {
	connPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return &DB{pool: connPool}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	err := db.pool.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
