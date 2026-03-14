package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/newtoallofthis123/moni/internal/models"
)

// TransactionInsert creates a transaction and updates the account balance atomically.
// For expenses, balance decreases; for income, balance increases.
func TransactionInsert(conn *sql.DB, accountID int64, categoryID *int64, txnType string, amount float64, note string) (models.Transaction, error) {
	var txn models.Transaction

	tx, err := conn.Begin()
	if err != nil {
		return txn, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = tx.QueryRow(
		`INSERT INTO transactions (account_id, category_id, type, amount, note)
		 VALUES (?, ?, ?, ?, ?)
		 RETURNING id, account_id, category_id, type, amount, note, date, created_at`,
		accountID, categoryID, txnType, amount, note,
	).Scan(&txn.ID, &txn.AccountID, &txn.CategoryID, &txn.Type, &txn.Amount, &txn.Note, &txn.Date, &txn.CreatedAt)
	if err != nil {
		return txn, fmt.Errorf("insert transaction: %w", err)
	}

	delta := amount
	if txnType == "expense" {
		delta = -amount
	}
	if err := AccountUpdateBalance(tx, accountID, delta); err != nil {
		return txn, err
	}

	if err := tx.Commit(); err != nil {
		return txn, fmt.Errorf("commit transaction: %w", err)
	}

	return txn, nil
}

// TransactionList returns transactions with optional filters.
// catName filters by category name (empty = no filter).
// since filters by date (zero time = no filter).
func TransactionList(conn *sql.DB, catName string, since time.Time) ([]models.Transaction, error) {
	query := `SELECT t.id, t.account_id, t.category_id, t.type, t.amount, t.note, t.date, t.created_at,
	                  a.name, COALESCE(c.name, '')
	           FROM transactions t
	           JOIN accounts a ON a.id = t.account_id
	           LEFT JOIN categories c ON c.id = t.category_id`

	var conditions []string
	var args []any

	if catName != "" {
		conditions = append(conditions, "c.name = ?")
		args = append(args, catName)
	}
	if !since.IsZero() {
		conditions = append(conditions, "t.date >= ?")
		args = append(args, since.Format("2006-01-02"))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.date DESC, t.id DESC"

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	var txns []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.CategoryID, &t.Type, &t.Amount, &t.Note, &t.Date, &t.CreatedAt, &t.AccountName, &t.CategoryName); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}
