package sql

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/jackc/pgx/v5"
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
	dir  fs.FS
	conn *pgx.Conn
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
		m, err := readMigration(m.dir, path)
		if err != nil {
			return fmt.Errorf("Unable to read migration at %s. %w", path, err)
		}
		files = append(files, *m)
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
	for _, migration := range files {
		err := m.executeMigration(ctx, migration)
		if err != nil {
			return fmt.Errorf("migrations failed: %w", err)
		}
	}

	return nil
}

func readMigration(dir fs.FS, path string) (*migration, error) {
	m := &migration{}
	version, err := strconv.Atoi(filepath.Base(path)[:4])
	if err != nil {
		return nil, fmt.Errorf("invalid version format: %w", err)
	}
	m.version = version
	contents, err := fs.ReadFile(dir, path)
	if err != nil {
		return nil, fmt.Errorf("could not read migration: %w", err)
	}
	m.sql = string(contents)
	return m, nil
}

func (svc *migrator) executeMigration(ctx context.Context, m migration) error {
	tx, err := svc.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not obtain transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = svc.conn.Exec(ctx, m.sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration %v: %w", m.version, err)
	}
	tx.Commit(ctx)
	return nil
}

type byVersion []migration

func (xs byVersion) Len() int           { return len(xs) }
func (xs byVersion) Swap(i, j int)      { xs[i], xs[j] = xs[j], xs[i] }
func (xs byVersion) Less(i, j int) bool { return xs[i].version < xs[j].version }
