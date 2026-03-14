package db

import (
	"testing"
)

func TestPersonInsertAndList(t *testing.T) {
	conn := testDB(t)

	p, err := PersonInsert(conn, "alice", "555-1234")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if p.Name != "alice" || p.Phone != "555-1234" {
		t.Errorf("unexpected person: %+v", p)
	}

	persons, _ := PersonList(conn)
	if len(persons) != 1 {
		t.Fatalf("expected 1 person, got %d", len(persons))
	}
}

func TestPersonDuplicateName(t *testing.T) {
	conn := testDB(t)

	PersonInsert(conn, "bob", "")
	_, err := PersonInsert(conn, "bob", "555-9999")
	if err == nil {
		t.Error("expected error for duplicate person name")
	}
}

func TestPersonGetByName(t *testing.T) {
	conn := testDB(t)

	PersonInsert(conn, "carol", "")
	p, err := PersonGetByName(conn, "carol")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if p.Name != "carol" {
		t.Errorf("expected 'carol', got %q", p.Name)
	}

	_, err = PersonGetByName(conn, "nobody")
	if err == nil {
		t.Error("expected error for nonexistent person")
	}
}

func TestPersonHistory(t *testing.T) {
	conn := testDB(t)

	p, _ := PersonInsert(conn, "dave", "")
	acct, _ := AccountInsert(conn, "main", "bank")

	// Add a transaction and link it
	txn, _ := TransactionInsert(conn, acct.ID, nil, "expense", 100, "dinner")
	TransactionPersonLink(conn, txn.ID, p.ID, "split bill")

	// Add a debt
	DebtInsert(conn, p.ID, 50, "they_owe", "cab")

	h, err := PersonHistory(conn, p.ID)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(h.Transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(h.Transactions))
	}
	if h.Transactions[0].LinkNote != "split bill" {
		t.Errorf("expected link note 'split bill', got %q", h.Transactions[0].LinkNote)
	}
	if len(h.Debts) != 1 {
		t.Errorf("expected 1 debt, got %d", len(h.Debts))
	}
}
