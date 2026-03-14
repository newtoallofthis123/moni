package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// testFileDB creates a real SQLite database in a temp directory with migrations applied.
// Returns the connection and the path to the .db file.
func testFileDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "moni.db")

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		t.Fatalf("wal mode: %v", err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys=ON"); err != nil {
		t.Fatalf("foreign keys: %v", err)
	}
	if err := Migrate(conn); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn, dbPath
}

func TestBackupCreateAndList(t *testing.T) {
	conn, dbPath := testFileDB(t)

	path, err := backupCreate(conn, dbPath)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("backup file not found: %v", err)
	}

	infos, err := backupList(backupDirFor(dbPath))
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(infos) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(infos))
	}
	if infos[0].Size == 0 {
		t.Error("backup size should be > 0")
	}
}

func TestBackupCreateOverwrite(t *testing.T) {
	conn, dbPath := testFileDB(t)

	if _, err := backupCreate(conn, dbPath); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := backupCreate(conn, dbPath); err != nil {
		t.Fatalf("second create: %v", err)
	}

	infos, err := backupList(backupDirFor(dbPath))
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(infos) != 1 {
		t.Errorf("expected 1 backup after overwrite, got %d", len(infos))
	}
}

func TestBackupDelete(t *testing.T) {
	conn, dbPath := testFileDB(t)
	dir := backupDirFor(dbPath)

	if _, err := backupCreate(conn, dbPath); err != nil {
		t.Fatalf("create: %v", err)
	}

	infos, _ := backupList(dir)
	if len(infos) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(infos))
	}

	if err := backupDelete(dir, infos[0].Date); err != nil {
		t.Fatalf("delete: %v", err)
	}

	infos, _ = backupList(dir)
	if len(infos) != 0 {
		t.Errorf("expected 0 backups after delete, got %d", len(infos))
	}

	// Delete nonexistent
	if err := backupDelete(dir, "1999-01-01"); err == nil {
		t.Error("expected error deleting nonexistent backup")
	}
}

func TestBackupPrune(t *testing.T) {
	conn, dbPath := testFileDB(t)
	dir := backupDirFor(dbPath)

	// Create the backup dir and fake multiple dated backups
	if _, err := backupCreate(conn, dbPath); err != nil {
		t.Fatalf("create: %v", err)
	}
	// Copy today's backup to simulate older backups
	todayInfos, _ := backupList(dir)
	src := todayInfos[0].Path
	for _, date := range []string{"2025-01-01", "2025-06-15", "2025-12-31"} {
		if err := copyFile(src, filepath.Join(dir, date+".db")); err != nil {
			t.Fatalf("copy fake backup: %v", err)
		}
	}

	infos, _ := backupList(dir)
	if len(infos) != 4 {
		t.Fatalf("expected 4 backups, got %d", len(infos))
	}

	removed, err := backupPrune(dir, 2)
	if err != nil {
		t.Fatalf("prune: %v", err)
	}
	if len(removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(removed))
	}

	infos, _ = backupList(dir)
	if len(infos) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(infos))
	}

	// The two newest should remain
	if infos[0].Date != "2025-12-31" {
		t.Errorf("expected oldest remaining 2025-12-31, got %s", infos[0].Date)
	}
}

func TestBackupPruneKeepZero(t *testing.T) {
	_, dbPath := testFileDB(t)
	dir := backupDirFor(dbPath)

	_, err := backupPrune(dir, 0)
	if err == nil {
		t.Error("expected error for keep=0")
	}
}

func TestBackupRestore(t *testing.T) {
	conn, dbPath := testFileDB(t)

	// Insert a bucket
	_, err := BucketInsert(conn, "savings", 1000)
	if err != nil {
		t.Fatalf("insert bucket: %v", err)
	}

	// Create backup
	if _, err := backupCreate(conn, dbPath); err != nil {
		t.Fatalf("create backup: %v", err)
	}

	// Delete the bucket
	buckets, _ := BucketList(conn)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if err := BucketDelete(conn, buckets[0].ID); err != nil {
		t.Fatalf("delete bucket: %v", err)
	}
	buckets, _ = BucketList(conn)
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets after delete, got %d", len(buckets))
	}

	// Get the backup date
	infos, _ := backupList(backupDirFor(dbPath))
	if len(infos) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(infos))
	}

	// Restore — this closes conn
	if err := backupRestore(conn, dbPath, infos[0].Date); err != nil {
		t.Fatalf("restore: %v", err)
	}

	// Reopen and verify bucket is back
	conn2, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer conn2.Close()

	buckets, _ = BucketList(conn2)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket after restore, got %d", len(buckets))
	}
	if buckets[0].Name != "savings" {
		t.Errorf("expected bucket name 'savings', got %q", buckets[0].Name)
	}
}
