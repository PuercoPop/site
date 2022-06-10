package main

import (
	"context"
	"log"

	"github.com/PuercoPop/site"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewAddUserCmd(dbpath string) *ffcli.Command {
	cmd := &ffcli.Command{
		Name:       "addadmin",
		ShortUsage: "swikictl addadmin",
		Exec: func(ctx context.Context, args []string) error {
			db, err := site.NewDB(ctx, dbpath)
			if err != nil {
				log.Fatalf("Could not connect to database. %s", err)
			}
			svc := site.NewUserService(db)
			email := "pirata@gmail.com"
			password := "yohoho"
			err = svc.CreateAdmin(ctx, email, password)
			return err

		},
	}
	return cmd
}
