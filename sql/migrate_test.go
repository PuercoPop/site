package sql

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
)

// Test scenarios
// 1. One migration creates a table

// 2. A failed migration, two tables, error in the second table
// statement. Assert that neither table is created. Checking that migrations are
// transnational.

// 3. Two migrations. First one succeeds. Second one fails. Checks that the
// transaction is per migration only.

// TODO(javier): Use podrick instead: https://github.com/uw-labs/podrick
func setupDB(t *testing.T) (*pgx.Conn, func()) {
	t.Helper()
	_, ok := os.LookupEnv("SLOW_TESTS")
	if !ok {
		// t.Skipf("%s requires a database to run. Set SLOW_TESTS to run this test.", t.Name())
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}
	// https://github.com/ory/dockertest/blob/v3/examples/PostgreSQL.md
	// Image information https://hub.docker.com/_/postgres
	resource, err := pool.Run("postgres", "15.1-alpine", []string{
		"POSTGRES_PASSWORD=P455W0RD",
		"POSTGRES_DB=migrate_test",
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	url := fmt.Sprintf("postgres://postgres:P455W0RD@%s/migrate_test?sslmode=disable", hostAndPort)
	err = resource.Expire(60)
	if err != nil {
		t.Fatalf("Could not set an expiration time to resource: %s", err)
	}
	pool.MaxWait = 120 * time.Second
	var conn *pgx.Conn
	if err = pool.Retry(func() error {
		conn, err = pgx.Connect(context.Background(), url)
		if err != nil {
			return err
		}
		return conn.Ping(context.Background())
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	}
	return conn, cleanup
}

func TestSortMigrations(t *testing.T) {
	migrations := byVersion{
		migration{version: 3},
		migration{version: 1},
		migration{version: 2},
	}
	sort.Sort(migrations)
	prevVersion := 0
	for _, m := range migrations {
		if prevVersion > m.version {
			t.Errorf("Migrations are not ordered. %v was before %v", prevVersion, m.version)
		}
		prevVersion = m.version
	}

}

func TestReadMigration(t *testing.T) {
	want := struct {
		version  int
		sql      string
		checksum string
	}{version: 1, sql: `CREATE TABLE users (
  user_id integer GENERATED always AS IDENTITY PRIMARY KEY
);
`, checksum: ""}
	m, err := readMigration(os.DirFS("./testdata/1/"), "0001_foo.sql")
	if err != nil {
		t.Fatalf("Could not read migration file. %s", err)
	}
	if m.version != want.version {
		t.Errorf("Version number mismatch (-want, +got): -%v, +%v", want.version, m.version)
	}
	if diff := cmp.Diff(want.sql, m.sql); diff != "" {
		t.Errorf("sql code mismatch (-want, +got): %s", diff)
	}
	// TODO(javier): Check the checksum matches
}

// Integration Tests. Rename TestRun or TestMigratorRun.
func TestMigrator(t *testing.T) {
	ctx := context.Background()
	t.Run("A single migration", func(t *testing.T) {
		conn, cleanup := setupDB(t)
		defer cleanup()
		m := migrator{conn: conn, dir: os.DirFS("./testdata/1")}
		err := conn.Ping(ctx)
		if err != nil {
			t.Fatalf("Could not connect to the database. %s", err)
		}
		// Run the migration
		if err := m.Run(ctx); err != nil {
			t.Fatalf("Failed to run the migration. %s", err)
		}
		// Check that the table exists

		rows, err := conn.Query(ctx, "SELECT table_name from information_schema.tables where table_schema = 'public'")
		if err != nil {
			t.Fatalf("Failed to check the database state. %s", err)
		}
		tables, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			t.Fatalf("Failed to collect rows. %s", err)
		}
		fmt.Printf("Tables: %v", tables)
		// select count(*) from information_schemata.tables where table_name = 'user' and schema = 'public'
		rows, err = conn.Query(ctx, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
		if err != nil {
			t.Fatalf("Could not collect table names. %s", err)
		}
		got, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			t.Fatalf("Failed to collect got. %s", err)
		}
		less := func(a, b string) bool { return a < b }
		if diff := cmp.Diff([]string{"users"}, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("Schema mismatch: (-want, +got): %s", diff)
		}
	})
}
