package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// RecurringInsert creates a new recurring item.
func RecurringInsert(db *sql.DB, description string, amount float64, categoryID *int64, frequency string, dueDay int, txnType string) (models.Recurring, error) {
	var r models.Recurring
	err := db.QueryRow(
		`INSERT INTO recurring (description, amount, category_id, frequency, due_day, type)
		 VALUES (?, ?, ?, ?, ?, ?)
		 RETURNING id, description, amount, category_id, frequency, due_day, type, active, created_at`,
		description, amount, categoryID, frequency, dueDay, txnType,
	).Scan(&r.ID, &r.Description, &r.Amount, &r.CategoryID, &r.Frequency, &r.DueDay, &r.Type, &r.Active, &r.CreatedAt)
	if err != nil {
		return r, fmt.Errorf("insert recurring: %w", err)
	}
	return r, nil
}

// RecurringDeactivate sets a recurring item as inactive.
func RecurringDeactivate(db *sql.DB, id int64) error {
	res, err := db.Exec(`UPDATE recurring SET active = 0 WHERE id = ? AND active = 1`, id)
	if err != nil {
		return fmt.Errorf("deactivate recurring: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("recurring item %d not found or already inactive", id)
	}
	return nil
}

// RecurringList returns all active recurring items with category names.
func RecurringList(db *sql.DB) ([]models.Recurring, error) {
	rows, err := db.Query(
		`SELECT r.id, r.description, r.amount, r.category_id, r.frequency, r.due_day, r.type, r.active, r.created_at,
		        COALESCE(c.name, '')
		 FROM recurring r
		 LEFT JOIN categories c ON c.id = r.category_id
		 WHERE r.active = 1
		 ORDER BY r.due_day, r.description`,
	)
	if err != nil {
		return nil, fmt.Errorf("list recurring: %w", err)
	}
	defer rows.Close()

	var items []models.Recurring
	for rows.Next() {
		var r models.Recurring
		if err := rows.Scan(&r.ID, &r.Description, &r.Amount, &r.CategoryID, &r.Frequency, &r.DueDay, &r.Type, &r.Active, &r.CreatedAt, &r.CategoryName); err != nil {
			return nil, fmt.Errorf("scan recurring: %w", err)
		}
		items = append(items, r)
	}
	return items, rows.Err()
}
