package sql

import (
	"context"
	"io/fs"
)

type migrator struct {
	// connection to database
	// directory where migrations reside
	migrationDir fs.FS
}

func (m *migrator) Run(ctx context.Context) {

}
