package site

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Could not parse the PG connection URL. %s", err)
	}
	db, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		log.Fatalf("%s", err)
	}
	// todo(javier): ping the db
	// func (p *Pool) Ping(ctx context.Context) error
	return db, nil
}
