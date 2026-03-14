package db

import (
	"testing"
)

func TestTransactionPersonLink(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "main", "bank")
	p, _ := PersonInsert(conn, "alice", "")
	txn, _ := TransactionInsert(conn, acct.ID, nil, "expense", 100, "dinner")

	tp, err := TransactionPersonLink(conn, txn.ID, p.ID, "split")
	if err != nil {
		t.Fatalf("link: %v", err)
	}
	if tp.TransactionID != txn.ID || tp.PersonID != p.ID || tp.Note != "split" {
		t.Errorf("unexpected link: %+v", tp)
	}

	// Duplicate link should fail (PK constraint)
	_, err = TransactionPersonLink(conn, txn.ID, p.ID, "again")
	if err == nil {
		t.Error("expected error for duplicate link")
	}
}

func TestTransactionDeleteCleansLinks(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "main", "bank")
	p, _ := PersonInsert(conn, "bob", "")
	txn, _ := TransactionInsert(conn, acct.ID, nil, "expense", 50, "")
	TransactionPersonLink(conn, txn.ID, p.ID, "")

	// Delete should clean up links too
	_, err := TransactionDelete(conn, txn.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Verify person history is empty
	h, _ := PersonHistory(conn, p.ID)
	if len(h.Transactions) != 0 {
		t.Errorf("expected 0 linked transactions after delete, got %d", len(h.Transactions))
	}
}
