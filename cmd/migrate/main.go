package main

import (
	"context"
	"flag"
	"log"

	"github.com/PuercoPop/site/sql"
	"github.com/jackc/pgx/v5"
)

var defaultdburl = "postures://swiki-dev"

func main() {
	ctx := context.Background()
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	dburl := fs.String("dburl", defaultdburl, "URL of the database to connect to")

	conn, err := pgx.Connect(ctx, *dburl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %s", err)
	}
	m := sql.NewMigrator(conn, sql.FSMigrations)
	if err := m.Run(ctx); err != nil {
		log.Fatalf("%s", err)
	}

}
