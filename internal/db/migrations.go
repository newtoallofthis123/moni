package db

import (
	"database/sql"
	"fmt"
)

const currentVersion = 1

// Migrate runs all pending migrations and seeds default data.
func Migrate(conn *sql.DB) error {
	var version int
	if err := conn.QueryRow("PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("cannot read schema version: %w", err)
	}

	if version >= currentVersion {
		return nil // already up to date
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback()

	if version < 1 {
		if err := migrateV1(tx); err != nil {
			return fmt.Errorf("migration v1 failed: %w", err)
		}
	}

	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", currentVersion)); err != nil {
		return fmt.Errorf("cannot update schema version: %w", err)
	}

	return tx.Commit()
}

func migrateV1(tx *sql.Tx) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL CHECK(type IN ('bank', 'cash', 'credit', 'wallet', 'other')),
			balance REAL NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL CHECK(type IN ('expense', 'income', 'both')),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER NOT NULL REFERENCES accounts(id),
			category_id INTEGER REFERENCES categories(id),
			type TEXT NOT NULL CHECK(type IN ('expense', 'income')),
			amount REAL NOT NULL CHECK(amount > 0),
			note TEXT,
			date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS recurring (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT NOT NULL,
			amount REAL NOT NULL CHECK(amount > 0),
			category_id INTEGER REFERENCES categories(id),
			frequency TEXT NOT NULL CHECK(frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
			due_day INTEGER NOT NULL,
			type TEXT NOT NULL CHECK(type IN ('expense', 'income')),
			active INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS buckets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			target REAL NOT NULL CHECK(target > 0),
			current REAL NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS persons (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			phone TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS debts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			person_id INTEGER NOT NULL REFERENCES persons(id),
			amount REAL NOT NULL CHECK(amount > 0),
			direction TEXT NOT NULL CHECK(direction IN ('i_owe', 'they_owe')),
			note TEXT,
			settled INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS transaction_persons (
			transaction_id INTEGER NOT NULL REFERENCES transactions(id),
			person_id INTEGER NOT NULL REFERENCES persons(id),
			note TEXT,
			PRIMARY KEY (transaction_id, person_id)
		)`,
	}

	for _, ddl := range tables {
		if _, err := tx.Exec(ddl); err != nil {
			return fmt.Errorf("table creation failed: %w", err)
		}
	}

	return seedCategories(tx)
}

func seedCategories(tx *sql.Tx) error {
	categories := []struct {
		name string
		typ  string
	}{
		// Expense
		{"food", "expense"},
		{"transport", "expense"},
		{"entertainment", "expense"},
		{"utilities", "expense"},
		{"rent", "expense"},
		{"health", "expense"},
		{"shopping", "expense"},
		{"subscriptions", "expense"},
		{"other", "expense"},
		// Income
		{"salary", "income"},
		{"freelance", "income"},
		{"gift", "income"},
		{"refund", "income"},
		{"other_income", "income"},
		// Both
		{"transfer", "both"},
	}

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO categories (name, type) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("cannot prepare category insert: %w", err)
	}
	defer stmt.Close()

	for _, c := range categories {
		if _, err := stmt.Exec(c.name, c.typ); err != nil {
			return fmt.Errorf("cannot seed category %q: %w", c.name, err)
		}
	}

	return nil
}
