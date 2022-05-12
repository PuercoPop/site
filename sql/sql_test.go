package sql

import (
	"testing"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
)

// Everything runs against an in-memory version.
// Test that the schema is valid.
// Test that created_at is initialized to the current timestamp

const dbopts sqlite.OpenFlags = sqlite.SQLITE_OPEN_MEMORY |
	sqlite.SQLITE_OPEN_READWRITE |
	sqlite.SQLITE_OPEN_CREATE |
	sqlite.SQLITE_OPEN_WAL |
	sqlite.SQLITE_OPEN_URI |
	sqlite.SQLITE_OPEN_NOMUTEX

func TestSchema(t *testing.T) {
	conn, err := sqlite.OpenConn(":memory:", dbopts)
	if err != nil {
		t.FailNow()
	}
	defer conn.Close()
	if err = sqlitex.ExecScript(conn, Schema); err != nil {
		t.Fatalf("Error executing the database schema.\n %s\n.", err)
	}
}
