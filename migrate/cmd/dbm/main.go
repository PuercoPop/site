package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"os"

	"github.com/PuercoPop/site/migrate"
	"github.com/jackc/pgx"
	"github.com/peterbourgon/ff"
)

//go:embed ../../migrations/*.sql
var FSMigrations embed.FS

func main() {
	ctx := context.Background()
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	var dburl = fs.String("d", "", "URL of the database to connect to")
	// var dir = fs.String("D", "", "The directory where the migrations are stored.")
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatalf("Could not parse flags: %s", err)
	}
	conn, err := pgx.Connect(ctx, *dburl)
	if err != nil {
		log.Fatalf("Could not connect to the database: %s", err)
	}
	m := migrate.New(conn, FSMigrations)
	if err := m.Setup(ctx); err != nil {
		log.Fatalf("%s", err)
	}
	if err := m.Run(ctx); err != nil {
		log.Fatalf("%s", err)
	}
}
