package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// TransactionPersonLink links a transaction to a person.
func TransactionPersonLink(db *sql.DB, transactionID, personID int64, note string) (models.TransactionPerson, error) {
	var tp models.TransactionPerson
	_, err := db.Exec(
		`INSERT INTO transaction_persons (transaction_id, person_id, note) VALUES (?, ?, ?)`,
		transactionID, personID, note,
	)
	if err != nil {
		return tp, fmt.Errorf("link transaction to person: %w", err)
	}
	tp.TransactionID = transactionID
	tp.PersonID = personID
	tp.Note = note
	return tp, nil
}

// TransactionExists checks whether a transaction with the given ID exists.
func TransactionExists(db *sql.DB, id int64) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM transactions WHERE id = ?)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check transaction: %w", err)
	}
	return exists, nil
}
