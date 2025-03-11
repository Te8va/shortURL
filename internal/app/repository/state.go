package repository

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var pg *sql.DB
var dsn string

func ConnectToPG(DSN string) error {
	var err error
	pg, err = sql.Open("pgx", DSN)
	if err != nil {
		return err
	}

	if err := pg.Ping(); err != nil {
		return err
	}

	_, err = pg.Exec("CREATE TABLE IF NOT EXISTS urlshrt(uuid SERIAL PRIMARY KEY, short_url TEXT, original_url TEXT)")
	if err != nil {
		return err
	}

	dsn = DSN
	return nil
}

func GetPgPtr() (*sql.DB, error) {
	if pg != nil {
		return pg, nil
	}
	return nil, errors.New("postgres was not initialized")
}
