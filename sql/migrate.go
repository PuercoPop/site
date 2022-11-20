package sql

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
)

type migration struct {
	version int
	// sql contains the sql code to run
	sql string
}

type migrator struct {
	// connection to database
	// directory where migrations reside
	dir fs.FS
}

func NewMigrator(migrationsDir fs.FS) *migrator {
	return &migrator{dir: migrationsDir}
}

var MigrationError = errors.New("migration error")

func (m *migrator) Run(ctx context.Context) error {
	// 1. List files in the dir, in lexicographical order.
	// 2. Execute each file.
	// 3. If there is an error running the file, return the error and exit.
	var files []migrations
	err := fs.WalkDir(m.dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		version, err := strconv.Atoi(filepath.Base(path)[:4])
		if err != nil {
			return fmt.Errorf("invalid version format: %w", err)
		}
		files = append(files, migration{version: version,
			file: path})
		return nil
	})
	// Sort the migrations
}
type byVersion []migration

func (xs byVersion) Len() int           { return len(xs) }
func (xs byVersion) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }
func (xs byVersion) Less(i, j int) bool { return xs[i].version < xs[j].version }
