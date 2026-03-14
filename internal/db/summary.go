package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/newtoallofthis123/moni/internal/models"
)

// MonthlySummary holds aggregated data for a month.
type MonthlySummary struct {
	Month         string                `json:"month"`
	TotalIncome   float64               `json:"total_income"`
	TotalExpenses float64               `json:"total_expenses"`
	Net           float64               `json:"net"`
	TopCategories []models.CategorySpend `json:"top_categories"`
}

// SummaryForMonth returns income/expense totals and top categories for a given month.
// month is the first day of the month.
func SummaryForMonth(db *sql.DB, month time.Time) (MonthlySummary, error) {
	start := month.Format("2006-01-02")
	end := month.AddDate(0, 1, 0).Format("2006-01-02")

	s := MonthlySummary{Month: month.Format("2006-01")}

	// Totals by type
	rows, err := db.Query(
		`SELECT type, COALESCE(SUM(amount), 0)
		 FROM transactions
		 WHERE date >= ? AND date < ?
		 GROUP BY type`, start, end,
	)
	if err != nil {
		return s, fmt.Errorf("summary totals: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var txnType string
		var total float64
		if err := rows.Scan(&txnType, &total); err != nil {
			return s, fmt.Errorf("scan summary: %w", err)
		}
		switch txnType {
		case "income":
			s.TotalIncome = total
		case "expense":
			s.TotalExpenses = total
		}
	}
	if err := rows.Err(); err != nil {
		return s, err
	}
	s.Net = s.TotalIncome - s.TotalExpenses

	// Top categories by spend
	catRows, err := db.Query(
		`SELECT COALESCE(c.name, 'uncategorized'), SUM(t.amount)
		 FROM transactions t
		 LEFT JOIN categories c ON c.id = t.category_id
		 WHERE t.type = 'expense' AND t.date >= ? AND t.date < ?
		 GROUP BY c.name
		 ORDER BY SUM(t.amount) DESC
		 LIMIT 5`, start, end,
	)
	if err != nil {
		return s, fmt.Errorf("summary categories: %w", err)
	}
	defer catRows.Close()

	for catRows.Next() {
		var cs models.CategorySpend
		if err := catRows.Scan(&cs.Category, &cs.Amount); err != nil {
			return s, fmt.Errorf("scan category spend: %w", err)
		}
		s.TopCategories = append(s.TopCategories, cs)
	}
	return s, catRows.Err()
}
