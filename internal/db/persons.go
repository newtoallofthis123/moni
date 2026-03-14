package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// PersonInsert creates a new person with the given name and optional phone.
func PersonInsert(db *sql.DB, name, phone string) (models.Person, error) {
	var p models.Person
	err := db.QueryRow(
		`INSERT INTO persons (name, phone) VALUES (?, ?) RETURNING id, name, phone, created_at`,
		name, phone,
	).Scan(&p.ID, &p.Name, &p.Phone, &p.CreatedAt)
	if err != nil {
		return p, fmt.Errorf("insert person: %w", err)
	}
	return p, nil
}

// PersonList returns all persons.
func PersonList(db *sql.DB) ([]models.Person, error) {
	rows, err := db.Query(`SELECT id, name, phone, created_at FROM persons ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list persons: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan person: %w", err)
		}
		persons = append(persons, p)
	}
	return persons, rows.Err()
}

// PersonGetByName returns the person with the given name.
func PersonGetByName(db *sql.DB, name string) (models.Person, error) {
	var p models.Person
	err := db.QueryRow(
		`SELECT id, name, phone, created_at FROM persons WHERE name = ?`, name,
	).Scan(&p.ID, &p.Name, &p.Phone, &p.CreatedAt)
	if err != nil {
		return p, fmt.Errorf("person %q not found: %w", name, err)
	}
	return p, nil
}

// PersonHistory returns a person's linked transactions and debts.
func PersonHistory(db *sql.DB, personID int64) (models.PersonHistory, error) {
	var h models.PersonHistory

	// Linked transactions
	rows, err := db.Query(
		`SELECT t.id, t.type, t.amount, t.note, t.date, COALESCE(c.name, ''), a.name, tp.note
		 FROM transaction_persons tp
		 JOIN transactions t ON t.id = tp.transaction_id
		 JOIN accounts a ON a.id = t.account_id
		 LEFT JOIN categories c ON c.id = t.category_id
		 WHERE tp.person_id = ?
		 ORDER BY t.date DESC`, personID,
	)
	if err != nil {
		return h, fmt.Errorf("list person transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pt models.PersonTransaction
		if err := rows.Scan(&pt.ID, &pt.Type, &pt.Amount, &pt.Note, &pt.Date, &pt.CategoryName, &pt.AccountName, &pt.LinkNote); err != nil {
			return h, fmt.Errorf("scan person transaction: %w", err)
		}
		h.Transactions = append(h.Transactions, pt)
	}
	if err := rows.Err(); err != nil {
		return h, err
	}

	// Debts
	debtRows, err := db.Query(
		`SELECT id, amount, direction, note, settled, created_at
		 FROM debts WHERE person_id = ?
		 ORDER BY created_at DESC`, personID,
	)
	if err != nil {
		return h, fmt.Errorf("list person debts: %w", err)
	}
	defer debtRows.Close()

	for debtRows.Next() {
		var d models.Debt
		if err := debtRows.Scan(&d.ID, &d.Amount, &d.Direction, &d.Note, &d.Settled, &d.CreatedAt); err != nil {
			return h, fmt.Errorf("scan person debt: %w", err)
		}
		h.Debts = append(h.Debts, d)
	}
	return h, debtRows.Err()
}
