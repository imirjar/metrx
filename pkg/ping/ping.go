package ping

import (
	"context"
	"log"
	"net"

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

func GetMyIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Print("client.go Dial ERROR", err)
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()
	return ip, nil
}
