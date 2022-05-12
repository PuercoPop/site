package main

import (
	"context"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func NewAddUserCmd() *ffcli.Command {
	cmd := &ffcli.Command{
		Name:       "addadmin",
		ShortUsage: "swikictl addadmin",
		Exec: func(ctx context.Context, args []string) error {
			return nil
		},
	}
	return cmd
}
