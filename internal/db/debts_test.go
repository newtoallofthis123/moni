package db

import (
	"testing"
)

func TestDebtInsertAndList(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "alice", "")
	d, err := DebtInsert(conn, p.ID, 100, "i_owe", "dinner")
	if err != nil {
		t.Fatalf("insert debt: %v", err)
	}
	if d.Amount != 100 || d.Direction != "i_owe" || d.Settled {
		t.Errorf("unexpected debt: %+v", d)
	}

	debts, err := DebtListOpen(conn)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(debts) != 1 {
		t.Fatalf("expected 1 debt, got %d", len(debts))
	}
	if debts[0].PersonName != "alice" {
		t.Errorf("expected person name 'alice', got %q", debts[0].PersonName)
	}
}

func TestDebtSettleFIFO(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "bob", "")

	// Add two debts: 50, then 100
	DebtInsert(conn, p.ID, 50, "they_owe", "coffee")
	DebtInsert(conn, p.ID, 100, "they_owe", "lunch")

	// Settle 75 — should fully settle first (50) and partially settle second (25 off 100)
	settled, err := DebtSettle(conn, p.ID, 75)
	if err != nil {
		t.Fatalf("settle: %v", err)
	}
	if settled != 75 {
		t.Errorf("expected 75 settled, got %.2f", settled)
	}

	debts, _ := DebtListOpen(conn)
	if len(debts) != 1 {
		t.Fatalf("expected 1 remaining debt, got %d", len(debts))
	}
	if debts[0].Amount != 75 { // 100 - 25 = 75
		t.Errorf("expected remaining 75, got %.2f", debts[0].Amount)
	}
}

func TestDebtSettleFullyClears(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "carol", "")
	DebtInsert(conn, p.ID, 100, "i_owe", "rent")

	settled, _ := DebtSettle(conn, p.ID, 100)
	if settled != 100 {
		t.Errorf("expected 100 settled, got %.2f", settled)
	}

	debts, _ := DebtListOpen(conn)
	if len(debts) != 0 {
		t.Errorf("expected 0 open debts, got %d", len(debts))
	}
}

func TestDebtSettleNoDebts(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "dave", "")
	_, err := DebtSettle(conn, p.ID, 50)
	if err == nil {
		t.Error("expected error settling when no debts exist")
	}
}

func TestDebtDelete(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "eve", "")
	d, _ := DebtInsert(conn, p.ID, 200, "they_owe", "")

	if err := DebtDelete(conn, d.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	debts, _ := DebtListOpen(conn)
	if len(debts) != 0 {
		t.Errorf("expected 0 debts after delete, got %d", len(debts))
	}

	// Delete nonexistent
	if err := DebtDelete(conn, 9999); err == nil {
		t.Error("expected error deleting nonexistent debt")
	}
}

func TestDebtInvalidDirection(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "frank", "")
	_, err := DebtInsert(conn, p.ID, 50, "invalid", "")
	if err == nil {
		t.Error("expected error for invalid direction")
	}
}
