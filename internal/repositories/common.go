package repositories

import "github.com/jackc/pgx/v5/pgconn"

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return true
	}
	return false
}
