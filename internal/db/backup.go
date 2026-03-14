package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// BackupInfo holds metadata about a single backup file.
type BackupInfo struct {
	Date string `json:"date"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// backupDirFor returns the backups directory derived from a DB path's parent.
func backupDirFor(dbPath string) string {
	return filepath.Join(filepath.Dir(dbPath), "backups")
}

// BackupDir returns the default backup directory (~/.moni/backups/) and ensures it exists.
func BackupDir() (string, error) {
	dbPath, err := DBPath()
	if err != nil {
		return "", err
	}
	dir := backupDirFor(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create backup directory: %w", err)
	}
	return dir, nil
}

// BackupCreate checkpoints WAL and copies the database to a date-stamped backup file.
func BackupCreate(conn *sql.DB) (string, error) {
	dbPath, err := DBPath()
	if err != nil {
		return "", err
	}
	return backupCreate(conn, dbPath)
}

func backupCreate(conn *sql.DB, dbPath string) (string, error) {
	if _, err := conn.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		return "", fmt.Errorf("wal checkpoint: %w", err)
	}

	dir := backupDirFor(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create backup directory: %w", err)
	}

	dest := filepath.Join(dir, time.Now().Format("2006-01-02")+".db")
	if err := copyFile(dbPath, dest); err != nil {
		return "", fmt.Errorf("backup copy: %w", err)
	}
	return dest, nil
}

// BackupList returns metadata for all backups in the default backup directory.
func BackupList() ([]BackupInfo, error) {
	dbPath, err := DBPath()
	if err != nil {
		return nil, err
	}
	return backupList(backupDirFor(dbPath))
}

func backupList(dir string) ([]BackupInfo, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.db"))
	if err != nil {
		return nil, err
	}

	sort.Strings(matches)
	infos := make([]BackupInfo, 0, len(matches))
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			continue
		}
		date := strings.TrimSuffix(filepath.Base(m), ".db")
		infos = append(infos, BackupInfo{
			Date: date,
			Path: m,
			Size: fi.Size(),
		})
	}
	return infos, nil
}

// BackupDelete removes a backup by date string (e.g. "2026-03-14").
func BackupDelete(date string) error {
	dbPath, err := DBPath()
	if err != nil {
		return err
	}
	return backupDelete(backupDirFor(dbPath), date)
}

func backupDelete(dir, date string) error {
	path := filepath.Join(dir, date+".db")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("backup %q not found", date)
	}
	return os.Remove(path)
}

// BackupPrune keeps the newest `keep` backups and removes the rest.
func BackupPrune(keep int) ([]string, error) {
	dbPath, err := DBPath()
	if err != nil {
		return nil, err
	}
	return backupPrune(backupDirFor(dbPath), keep)
}

func backupPrune(dir string, keep int) ([]string, error) {
	if keep <= 0 {
		return nil, fmt.Errorf("--keep must be a positive number")
	}

	infos, err := backupList(dir)
	if err != nil {
		return nil, err
	}

	if len(infos) <= keep {
		return nil, nil
	}

	// infos is sorted ascending by date; remove the oldest ones
	toRemove := infos[:len(infos)-keep]
	removed := make([]string, 0, len(toRemove))
	for _, info := range toRemove {
		if err := os.Remove(info.Path); err != nil {
			return removed, fmt.Errorf("removing %s: %w", info.Date, err)
		}
		removed = append(removed, info.Date)
	}
	return removed, nil
}

// BackupRestore replaces the main database with the specified backup.
// The caller must not use conn after this call (it will be closed).
func BackupRestore(conn *sql.DB, date string) error {
	dbPath, err := DBPath()
	if err != nil {
		return err
	}
	return backupRestore(conn, dbPath, date)
}

func backupRestore(conn *sql.DB, dbPath, date string) error {
	dir := backupDirFor(dbPath)
	src := filepath.Join(dir, date+".db")
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("backup %q not found", date)
	}

	if _, err := conn.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		return fmt.Errorf("wal checkpoint: %w", err)
	}
	conn.Close()

	if err := copyFile(src, dbPath); err != nil {
		return fmt.Errorf("restore copy: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
