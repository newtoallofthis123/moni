package db

import (
	"database/sql"
	"fmt"

	"github.com/newtoallofthis123/moni/internal/models"
)

// BucketInsert creates a new savings bucket.
func BucketInsert(db *sql.DB, name string, target float64) (models.Bucket, error) {
	var b models.Bucket
	err := db.QueryRow(
		`INSERT INTO buckets (name, target)
		 VALUES (?, ?)
		 RETURNING id, name, target, current, created_at`,
		name, target,
	).Scan(&b.ID, &b.Name, &b.Target, &b.Current, &b.CreatedAt)
	if err != nil {
		return b, fmt.Errorf("insert bucket: %w", err)
	}
	return b, nil
}

// BucketAddFunds adds an amount to a bucket's current balance.
func BucketAddFunds(db *sql.DB, bucketID int64, amount float64) (models.Bucket, error) {
	var b models.Bucket
	err := db.QueryRow(
		`UPDATE buckets SET current = current + ?
		 WHERE id = ?
		 RETURNING id, name, target, current, created_at`,
		amount, bucketID,
	).Scan(&b.ID, &b.Name, &b.Target, &b.Current, &b.CreatedAt)
	if err != nil {
		return b, fmt.Errorf("add funds to bucket: %w", err)
	}
	return b, nil
}

// BucketGetByName returns a bucket by name.
func BucketGetByName(db *sql.DB, name string) (models.Bucket, error) {
	var b models.Bucket
	err := db.QueryRow(
		`SELECT id, name, target, current, created_at FROM buckets WHERE name = ?`, name,
	).Scan(&b.ID, &b.Name, &b.Target, &b.Current, &b.CreatedAt)
	if err != nil {
		return b, fmt.Errorf("bucket %q not found: %w", name, err)
	}
	return b, nil
}

// BucketList returns all buckets.
func BucketList(db *sql.DB) ([]models.Bucket, error) {
	rows, err := db.Query(
		`SELECT id, name, target, current, created_at FROM buckets ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("list buckets: %w", err)
	}
	defer rows.Close()

	var buckets []models.Bucket
	for rows.Next() {
		var b models.Bucket
		if err := rows.Scan(&b.ID, &b.Name, &b.Target, &b.Current, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan bucket: %w", err)
		}
		buckets = append(buckets, b)
	}
	return buckets, rows.Err()
}
