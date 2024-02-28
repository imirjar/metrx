package ping

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingPgx(ctx context.Context, conn string) error {
	if conn == "" {
		log.Println(errConnectionString)
		return errConnectionString
	}

	db, err := sql.Open("pgx", conn)
	if err != nil {
		log.Print(errConnectionParams)
		return errConnectionParams
	}
	defer db.Close()

	return db.PingContext(ctx)
}
