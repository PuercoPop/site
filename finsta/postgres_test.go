package finsta

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func cleardb(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	_, err := db.Exec(ctx, "TRUNCATE sessions, users")
	if err != nil {
		t.Fatalf("Unable to restart database. %s", err)
	}

}

func requiresdb(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.TODO()
	// todo(javier): Get database name from environment variable
	pool, err := pgxpool.Connect(ctx, "postgres:///swiki-test")
	if err != nil {
		t.Fatalf("Could not connect to database %s", err)
	}
	return pool
}
