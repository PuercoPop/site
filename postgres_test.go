package swiki

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func setuptestdb(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	pool, err := pgxpool.Connect(context.TODO(), "postgres:///swiki-test")
	if err != nil {
		t.Fatalf("Could not connect to database %s", err)
	}
	return pool, func() { pool.Close() }
}
