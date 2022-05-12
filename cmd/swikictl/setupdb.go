package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"github.com/PuercoPop/swiki/sql"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func setupdb(dbpath string) error {
	fmt.Println("Calling setuppath with ", dbpath)
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
				log.Fatalf("Could not delete #{dbpath}", dbpath, err)
			}
		}
	}
	conn, err := sqlite.OpenConn(dbpath, 0)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = sqlitex.ExecScript(conn, sql.Schema); err != nil {
		log.Fatalf("Error executing the database schema. %s\n.", err)
	}
	return nil
}

func NewSetupDBCmd(dbpath string) *ffcli.Command {
	cmd := &ffcli.Command{
		Name:       "setupdb",
		ShortUsage: "swikictl setupdb",
		Exec: func(ctx context.Context, args []string) error {
			return setupdb(dbpath)
		},
	}
	return cmd
}
