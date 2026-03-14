package db

import (
	"testing"
	"time"
)

func TestTransactionInsertUpdatesBalance(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "checking", "bank")
	cat, _ := CategoryGetByName(conn, "salary")
	catID := cat.ID

	// Income increases balance
	_, err := TransactionInsert(conn, acct.ID, &catID, "income", 1000, "paycheck")
	if err != nil {
		t.Fatalf("insert income: %v", err)
	}
	updated, _ := AccountGetByName(conn, "checking")
	if updated.Balance != 1000 {
		t.Errorf("expected balance 1000, got %.2f", updated.Balance)
	}

	// Expense decreases balance
	food, _ := CategoryGetByName(conn, "food")
	foodID := food.ID
	_, err = TransactionInsert(conn, acct.ID, &foodID, "expense", 250, "groceries")
	if err != nil {
		t.Fatalf("insert expense: %v", err)
	}
	updated, _ = AccountGetByName(conn, "checking")
	if updated.Balance != 750 {
		t.Errorf("expected balance 750, got %.2f", updated.Balance)
	}
}

func TestTransactionDelete(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "wallet", "cash")
	cat, _ := CategoryGetByName(conn, "food")
	catID := cat.ID

	// Add income then expense
	TransactionInsert(conn, acct.ID, &catID, "income", 500, "")
	txn, _ := TransactionInsert(conn, acct.ID, &catID, "expense", 200, "lunch")

	// Balance should be 300
	a, _ := AccountGetByName(conn, "wallet")
	if a.Balance != 300 {
		t.Fatalf("expected 300, got %.2f", a.Balance)
	}

	// Delete expense should restore balance
	_, err := TransactionDelete(conn, txn.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	a, _ = AccountGetByName(conn, "wallet")
	if a.Balance != 500 {
		t.Errorf("expected 500 after delete, got %.2f", a.Balance)
	}

	// Delete nonexistent
	_, err = TransactionDelete(conn, 9999)
	if err == nil {
		t.Error("expected error deleting nonexistent transaction")
	}
}

func TestTransactionListFilters(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "main", "bank")
	food, _ := CategoryGetByName(conn, "food")
	salary, _ := CategoryGetByName(conn, "salary")
	foodID := food.ID
	salaryID := salary.ID

	TransactionInsert(conn, acct.ID, &foodID, "expense", 50, "lunch")
	TransactionInsert(conn, acct.ID, &salaryID, "income", 2000, "pay")
	TransactionInsert(conn, acct.ID, &foodID, "expense", 30, "dinner")

	// No filter — all 3
	all, _ := TransactionList(conn, "", time.Time{})
	if len(all) != 3 {
		t.Errorf("expected 3 transactions, got %d", len(all))
	}

	// Filter by category
	foodTxns, _ := TransactionList(conn, "food", time.Time{})
	if len(foodTxns) != 2 {
		t.Errorf("expected 2 food transactions, got %d", len(foodTxns))
	}

	// Filter by date (future date should return 0)
	future := time.Now().Add(24 * time.Hour)
	none, _ := TransactionList(conn, "", future)
	if len(none) != 0 {
		t.Errorf("expected 0 future transactions, got %d", len(none))
	}
}

func TestTransactionExists(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "test", "bank")
	txn, _ := TransactionInsert(conn, acct.ID, nil, "income", 100, "")

	exists, _ := TransactionExists(conn, txn.ID)
	if !exists {
		t.Error("expected transaction to exist")
	}

	exists, _ = TransactionExists(conn, 9999)
	if exists {
		t.Error("expected nonexistent transaction to not exist")
	}
}
