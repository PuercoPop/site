package sql

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type migration struct {
	version int
	// sql contains the sql code to run
	sql      string
	checksum string
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
	var files byVersion = make([]migration, 0)
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
			sql: path})
		return nil
	})
	if err != nil {
		return fmt.Errorf("Could not load migrations %w", err)
	}
	// Sort the migrations
	sort.Sort(files)
	fmt.Printf("migrations: %v\n", files)
	// Check if any migrations have been executed or tampered with
	// Execute migrations
	return nil
}

func readMigration(path string) (*migration, error) {
	m := &migration{}
	version, err := strconv.Atoi(filepath.Base(path)[:4])
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %w", err)
	}
	m.version = version
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read migration: %w", err)
	}
	m.sql = string(contents)
	return m, nil
}

type byVersion []migration

func (xs byVersion) Len() int           { return len(xs) }
func (xs byVersion) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }
func (xs byVersion) Less(i, j int) bool { return xs[i].version < xs[j].version }
