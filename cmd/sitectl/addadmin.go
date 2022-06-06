package main

import (
	"context"
	"log"

	"github.com/PuercoPop/site"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/crypto/bcrypt"
)

func NewAddUserCmd(dbpath string) *ffcli.Command {
	cmd := &ffcli.Command{
		Name:       "addadmin",
		ShortUsage: "swikictl addadmin",
		Exec: func(ctx context.Context, args []string) error {
			db, err := site.NewDB(ctx, dbpath)
			email := "pirata@gmail.com"
			password := "yohoho"
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Could not hash password. %s", err)
			}
			_, err = db.Exec(ctx, "INSERT INTO users (email, password, admin) VALUES ($1, $2, true)", email, hash)
			if err != nil {
				log.Fatalf("Could not insert record. %s", err)
			}
			return nil

		},
	}
	return cmd
}
