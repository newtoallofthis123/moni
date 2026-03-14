package db

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// testDB creates an in-memory SQLite database with all migrations applied.
func testDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys=ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	if err := Migrate(conn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestMigrate(t *testing.T) {
	conn := testDB(t)

	// Verify tables exist
	tables := []string{"accounts", "categories", "transactions", "recurring", "buckets", "persons", "debts", "transaction_persons"}
	for _, tbl := range tables {
		var name string
		err := conn.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tbl).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", tbl, err)
		}
	}

	// Verify seeded categories exist
	var count int
	conn.QueryRow(`SELECT COUNT(*) FROM categories`).Scan(&count)
	if count != 15 {
		t.Errorf("expected 15 seeded categories, got %d", count)
	}

	// Verify idempotent
	if err := Migrate(conn); err != nil {
		t.Errorf("second migration should be no-op: %v", err)
	}
}
