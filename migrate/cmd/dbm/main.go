package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"os"

	"github.com/PuercoPop/site/migrate"
	"github.com/jackc/pgx/v5"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	flagset := flag.NewFlagSet("migrate", flag.ExitOnError)
	var dburl = flagset.String("d", "", "URL of the database to connect to")
	var dir = flagset.String("D", "", "The directory where the migrations are stored.")
	err := ff.Parse(flagset, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatalf("Could not parse flags: %s", err)
	}
	var FSDir fs.FS
	if *dir == "" {
		FSDir = migrate.FSMigrations
	} else {
		FSDir = os.DirFS(*dir)
	}
	if *dburl == "" {
		log.Fatalf("Database URL must be provided.")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, *dburl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %s", err)
	}
	m := migrate.New(conn, FSDir)
	if err := m.Setup(ctx); err != nil {
		log.Fatalf("%s", err)
	}
	if err := m.Run(ctx); err != nil {
		log.Fatalf("%s", err)
	}
}
