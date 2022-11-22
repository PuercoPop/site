package sql

import (
	"context"
	_ "embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var Schema string

func NewDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Could not parse the PG connection URL. %s", err)
	}
	db, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		log.Fatalf("%s", err)
	}
	// TODO(javier): ping the db
	// func (p *Pool) Ping(ctx context.Context) error
	return db, nil
}
