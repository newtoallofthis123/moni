package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// DebtInsert creates a new debt record for a person.
func DebtInsert(db *sql.DB, personID int64, amount float64, direction, note string) (models.Debt, error) {
	var d models.Debt
	err := db.QueryRow(
		`INSERT INTO debts (person_id, amount, direction, note)
		 VALUES (?, ?, ?, ?)
		 RETURNING id, person_id, amount, direction, note, settled, created_at`,
		personID, amount, direction, note,
	).Scan(&d.ID, &d.PersonID, &d.Amount, &d.Direction, &d.Note, &d.Settled, &d.CreatedAt)
	if err != nil {
		return d, fmt.Errorf("insert debt: %w", err)
	}
	return d, nil
}

// DebtSettle settles debts for a person, oldest first (FIFO).
// Returns total amount actually settled.
func DebtSettle(conn *sql.DB, personID int64, amount float64) (float64, error) {
	tx, err := conn.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get unsettled debts for this person, oldest first
	rows, err := tx.Query(
		`SELECT id, amount FROM debts
		 WHERE person_id = ? AND settled = 0
		 ORDER BY created_at ASC`, personID,
	)
	if err != nil {
		return 0, fmt.Errorf("list unsettled debts: %w", err)
	}

	type debtRow struct {
		id     int64
		amount float64
	}
	var debts []debtRow
	for rows.Next() {
		var d debtRow
		if err := rows.Scan(&d.id, &d.amount); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan debt: %w", err)
		}
		debts = append(debts, d)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	if len(debts) == 0 {
		return 0, fmt.Errorf("no unsettled debts for this person")
	}

	remaining := amount
	settled := 0.0

	for _, d := range debts {
		if remaining <= 0 {
			break
		}
		if remaining >= d.amount {
			// Fully settle this debt
			_, err := tx.Exec(`UPDATE debts SET settled = 1 WHERE id = ?`, d.id)
			if err != nil {
				return 0, fmt.Errorf("settle debt: %w", err)
			}
			remaining -= d.amount
			settled += d.amount
		} else {
			// Partially settle: reduce amount, keep unsettled
			_, err := tx.Exec(`UPDATE debts SET amount = amount - ? WHERE id = ?`, remaining, d.id)
			if err != nil {
				return 0, fmt.Errorf("partial settle debt: %w", err)
			}
			settled += remaining
			remaining = 0
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return settled, nil
}

// DebtListOpen returns all unsettled debts with person names.
func DebtListOpen(conn *sql.DB) ([]models.Debt, error) {
	rows, err := conn.Query(
		`SELECT d.id, d.person_id, d.amount, d.direction, d.note, d.settled, d.created_at, p.name
		 FROM debts d
		 JOIN persons p ON p.id = d.person_id
		 WHERE d.settled = 0
		 ORDER BY p.name, d.created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list open debts: %w", err)
	}
	defer rows.Close()

	var debts []models.Debt
	for rows.Next() {
		var d models.Debt
		if err := rows.Scan(&d.ID, &d.PersonID, &d.Amount, &d.Direction, &d.Note, &d.Settled, &d.CreatedAt, &d.PersonName); err != nil {
			return nil, fmt.Errorf("scan debt: %w", err)
		}
		debts = append(debts, d)
	}
	return debts, rows.Err()
}
