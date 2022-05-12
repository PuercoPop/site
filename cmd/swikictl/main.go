package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func dbpath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".swiki.sqlite3")
}

func main() {
	fs := flag.NewFlagSet("swikictl", flag.ExitOnError)
	var dbpath = fs.String("d", dbpath(), "The path to the sqlite3 database")
	root := &ffcli.Command{
		ShortUsage:  "swikictl [-d database] <command>",
		ShortHelp:   "Control swiki",
		FlagSet:     fs,
		Subcommands: []*ffcli.Command{NewSetupDBCmd(*dbpath), NewAddUserCmd()},
	}
	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
