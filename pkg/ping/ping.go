package ping

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingPgx(ctx context.Context, conn string) bool {
	context := context.WithoutCancel(ctx)
	if conn == "" {
		log.Println("1 CONNECTION ISNT VALID")
		return false
	}

	db, err := sql.Open("pgx", conn)
	if err != nil {
		log.Println("2 CONNECTION ISNT VALID")
		return false
	}
	defer db.Close()

	err = db.PingContext(context)
	if err == nil {
		log.Println("3 CONNECTION IS VALID")
	} else {
		log.Println("4 CONNECTION ISNT VALID")
	}

	return err == nil
}
