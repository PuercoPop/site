package sql

import "io/fs"

type migrator struct {
	// connection to database
	// directory where migrations reside
	migrationDir fs.FS
}
