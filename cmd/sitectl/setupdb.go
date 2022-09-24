package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuercoPop/site/sql"
	"github.com/jackc/pgx/v4"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func setupdb(ctx context.Context, dbpath string) error {
	if _, err := os.Stat(dbpath); !os.IsNotExist(err) {
		fmt.Println("A previous version of the database was found. Delete it? [y/N]")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			fmt.Printf("Deleting %s\n", dbpath)
			err := os.Remove(dbpath)
			if err != nil {
				log.Fatalf("Could not delete %s. %s", dbpath, err)
			}
		}
	}
	conf, err := pgx.ParseConfig(dbpath)
	if err != nil {
		return err
	}
	conn, err := pgx.ConnectConfig(ctx, conf)
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, sql.Schema)
	if err != nil {
		log.Fatalf("Error executing the database schema. %s\n.", err)
	}
	return nil
}

func NewSetupDBCmd(dbpath string) *ffcli.Command {
	cmd := &ffcli.Command{
		Name:       "setupdb",
		ShortUsage: "swikictl setupdb",
		Exec: func(ctx context.Context, args []string) error {
			return setupdb(ctx, dbpath)
		},
	}
	return cmd
}
