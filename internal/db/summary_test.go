package db

import (
	"testing"
	"time"
)

func TestSummaryForMonth(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "main", "bank")
	food, _ := CategoryGetByName(conn, "food")
	salary, _ := CategoryGetByName(conn, "salary")
	transport, _ := CategoryGetByName(conn, "transport")
	foodID := food.ID
	salaryID := salary.ID
	transportID := transport.ID

	// Add some transactions — these have today's date
	TransactionInsert(conn, acct.ID, &salaryID, "income", 5000, "pay", time.Time{})
	TransactionInsert(conn, acct.ID, &foodID, "expense", 200, "groceries", time.Time{})
	TransactionInsert(conn, acct.ID, &foodID, "expense", 100, "restaurant", time.Time{})
	TransactionInsert(conn, acct.ID, &transportID, "expense", 50, "bus", time.Time{})

	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	s, err := SummaryForMonth(conn, month)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}

	if s.TotalIncome != 5000 {
		t.Errorf("expected income 5000, got %.2f", s.TotalIncome)
	}
	if s.TotalExpenses != 350 {
		t.Errorf("expected expenses 350, got %.2f", s.TotalExpenses)
	}
	if s.Net != 4650 {
		t.Errorf("expected net 4650, got %.2f", s.Net)
	}
	if len(s.TopCategories) != 2 {
		t.Errorf("expected 2 expense categories, got %d", len(s.TopCategories))
	}
	// Food should be first (300 > 50)
	if len(s.TopCategories) > 0 && s.TopCategories[0].Category != "food" {
		t.Errorf("expected top category 'food', got %q", s.TopCategories[0].Category)
	}
}

func TestSummaryEmptyMonth(t *testing.T) {
	conn := testDB(t)

	// No transactions — summary should return zeros
	month := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	s, err := SummaryForMonth(conn, month)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if s.TotalIncome != 0 || s.TotalExpenses != 0 || s.Net != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
	if len(s.TopCategories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(s.TopCategories))
	}
}
