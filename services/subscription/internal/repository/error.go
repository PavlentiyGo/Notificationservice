package repository

import (
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsPgErr(err error, pgCode string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgCode {
			return true
		}
	}
	log.Println(1)
	return false
}
