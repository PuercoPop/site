package main

import (
	"context"
	"log"

	"github.com/PuercoPop/site/finsta"
	"github.com/PuercoPop/site/sql"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewAddUserCmd(dbpath string) *ffcli.Command {
	cmd := &ffcli.Command{
		Name:       "addadmin",
		ShortUsage: "swikictl addadmin",
		Exec: func(ctx context.Context, args []string) error {
			db, err := sql.NewDB(ctx, dbpath)
			if err != nil {
				log.Fatalf("Could not connect to database. %s", err)
			}
			svc := finsta.NewUserService(db)
			email := "pirata@gmail.com"
			password := "yohoho"
			err = svc.CreateAdmin(ctx, email, password)
			return err

		},
	}
	return cmd
}
