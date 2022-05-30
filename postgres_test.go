package swiki

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func setuptestdb(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	ctx := context.TODO()
	pool, err := pgxpool.Connect(ctx, "postgres:///swiki-test")
	if err != nil {
		t.Fatalf("Could not connect to database %s", err)
	}
	tx, err := pool.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		t.Fatalf("Failed to acquire the transaction. %s", err)
	}
	return pool, func() {
		tx.Rollback(ctx)
		pool.Close()
	}
}
