package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// CategoryInsert creates a new category with the given name and type.
func CategoryInsert(db *sql.DB, name, catType string) (models.Category, error) {
	var cat models.Category
	err := db.QueryRow(
		`INSERT INTO categories (name, type) VALUES (?, ?) RETURNING id, name, type, created_at`,
		name, catType,
	).Scan(&cat.ID, &cat.Name, &cat.Type, &cat.CreatedAt)
	if err != nil {
		return cat, fmt.Errorf("insert category: %w", err)
	}
	return cat, nil
}

// CategoryList returns all categories.
func CategoryList(db *sql.DB) ([]models.Category, error) {
	rows, err := db.Query(`SELECT id, name, type, created_at FROM categories ORDER BY type, name`)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

// CategoryGetByName returns the category with the given name.
func CategoryGetByName(db *sql.DB, name string) (models.Category, error) {
	var c models.Category
	err := db.QueryRow(
		`SELECT id, name, type, created_at FROM categories WHERE name = ?`, name,
	).Scan(&c.ID, &c.Name, &c.Type, &c.CreatedAt)
	if err != nil {
		return c, fmt.Errorf("get category %q: %w", name, err)
	}
	return c, nil
}
