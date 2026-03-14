package db

import (
	"testing"
)

func TestRecurringInsertAndList(t *testing.T) {
	conn := testDB(t)

	cat, _ := CategoryGetByName(conn, "rent")
	catID := cat.ID

	r, err := RecurringInsert(conn, "Monthly rent", 1200, &catID, "monthly", 1, "expense")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if r.Description != "Monthly rent" || r.Amount != 1200 || !r.Active {
		t.Errorf("unexpected recurring: %+v", r)
	}

	items, _ := RecurringList(conn)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].CategoryName != "rent" {
		t.Errorf("expected category 'rent', got %q", items[0].CategoryName)
	}
}

func TestRecurringDeactivate(t *testing.T) {
	conn := testDB(t)

	r, _ := RecurringInsert(conn, "Netflix", 15, nil, "monthly", 15, "expense")

	if err := RecurringDeactivate(conn, r.ID); err != nil {
		t.Fatalf("deactivate: %v", err)
	}

	// Should not appear in active list
	items, _ := RecurringList(conn)
	if len(items) != 0 {
		t.Errorf("expected 0 active items, got %d", len(items))
	}

	// Deactivating again should error
	if err := RecurringDeactivate(conn, r.ID); err == nil {
		t.Error("expected error deactivating already inactive item")
	}
}

func TestRecurringInvalidFrequency(t *testing.T) {
	conn := testDB(t)

	_, err := RecurringInsert(conn, "bad", 10, nil, "biweekly", 1, "expense")
	if err == nil {
		t.Error("expected error for invalid frequency")
	}
}
