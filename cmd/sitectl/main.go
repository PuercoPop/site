package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var defaultdburl = "postgres:///swiki-dev"

func main() {
	var fs = flag.NewFlagSet("swikictl", flag.ExitOnError)
	var dbpath = fs.String("d", defaultdburl, "The path to the sqlite3 database")
	root := &ffcli.Command{
		ShortUsage:  "swikictl [-d database] <command>",
		ShortHelp:   "Control swiki",
		FlagSet:     fs,
		Subcommands: []*ffcli.Command{NewSetupDBCmd(*dbpath), NewAddUserCmd(*dbpath)},
	}
	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
