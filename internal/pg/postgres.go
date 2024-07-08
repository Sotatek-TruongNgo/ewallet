package pg

import (
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(
	host string,
	database string,
	username string,
	password string,
) (*sqlx.DB, error) {
	uri := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		url.QueryEscape(username),
		url.QueryEscape(password),
		host,
		database,
	)

	config, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("parse postgres connection string failed: %w", err)
	}

	db := sqlx.NewDb(stdlib.OpenDB(*config), "pgx")
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping to database failed: %w", err)
	}
	db.SetMaxOpenConns(5000)
	db.SetMaxIdleConns(1000)

	return db, nil
}
