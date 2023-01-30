package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/PuercoPop/site/migrate"
	"github.com/jackc/pgx/v5"
	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	var dburl = fs.String("d", "", "URL of the database to connect to")
	// var dir = fs.String("D", "", "The directory where the migrations are stored.")
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatalf("Could not parse flags: %s", err)
	}
	dir := migrate.FSMigrations
	// if *dir == "" {
	// 	dir = migrate.FSMigrations
	// }
	if *dburl == "" {
		log.Fatalf("Database URL must be provided.")
	}
	ctx := context.Background()
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
