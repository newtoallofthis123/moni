package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// AccountInsert creates a new account with the given name and type.
func AccountInsert(db *sql.DB, name, acctType string) (models.Account, error) {
	var acct models.Account
	err := db.QueryRow(
		`INSERT INTO accounts (name, type) VALUES (?, ?) RETURNING id, name, type, balance, created_at`,
		name, acctType,
	).Scan(&acct.ID, &acct.Name, &acct.Type, &acct.Balance, &acct.CreatedAt)
	if err != nil {
		return acct, fmt.Errorf("insert account: %w", err)
	}
	return acct, nil
}

// AccountList returns all accounts.
func AccountList(db *sql.DB) ([]models.Account, error) {
	rows, err := db.Query(`SELECT id, name, type, balance, created_at FROM accounts ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

// AccountGetByName returns the account with the given name.
func AccountGetByName(db *sql.DB, name string) (models.Account, error) {
	var a models.Account
	err := db.QueryRow(
		`SELECT id, name, type, balance, created_at FROM accounts WHERE name = ?`, name,
	).Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.CreatedAt)
	if err != nil {
		return a, fmt.Errorf("get account %q: %w", name, err)
	}
	return a, nil
}

// AccountGetFirst returns the first account (by ID), used as default.
func AccountGetFirst(db *sql.DB) (models.Account, error) {
	var a models.Account
	err := db.QueryRow(
		`SELECT id, name, type, balance, created_at FROM accounts ORDER BY id LIMIT 1`,
	).Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.CreatedAt)
	if err != nil {
		return a, fmt.Errorf("no accounts found: %w", err)
	}
	return a, nil
}

// AccountUpdateBalance adjusts the balance of an account by delta (positive or negative).
func AccountUpdateBalance(tx *sql.Tx, accountID int64, delta float64) error {
	_, err := tx.Exec(`UPDATE accounts SET balance = balance + ? WHERE id = ?`, delta, accountID)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	return nil
}
