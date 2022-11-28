package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/PuercoPop/site/sql"
	"github.com/jackc/pgx/v5"
	"github.com/peterbourgon/ff/v3"
)

var defaultdburl = "postgres://swiki-dev"

func main() {
	ctx := context.Background()
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	var dburl = fs.String("dburl", defaultdburl, "URL of the database to connect to")
	err := ff.Parse(fs, os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse flags: %s", err)
	}
	conn, err := pgx.Connect(ctx, *dburl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %s", err)
	}
	m := sql.NewMigrator(conn, sql.FSMigrations)
	if err := m.Setup(ctx); err != nil {
		log.Fatalf("%s", err)
	}
	if err := m.Run(ctx); err != nil {
		log.Fatalf("%s", err)
	}
}
