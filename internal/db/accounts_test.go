package db

import (
	"testing"
)

func TestAccountInsertAndList(t *testing.T) {
	conn := testDB(t)

	acct, err := AccountInsert(conn, "checking", "bank")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if acct.Name != "checking" || acct.Type != "bank" || acct.Balance != 0 {
		t.Errorf("unexpected account: %+v", acct)
	}

	accounts, err := AccountList(conn)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}
	if accounts[0].ID != acct.ID {
		t.Errorf("expected id %d, got %d", acct.ID, accounts[0].ID)
	}
}

func TestAccountDuplicateName(t *testing.T) {
	conn := testDB(t)

	if _, err := AccountInsert(conn, "cash", "cash"); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if _, err := AccountInsert(conn, "cash", "wallet"); err == nil {
		t.Error("expected error for duplicate name")
	}
}

func TestAccountGetByName(t *testing.T) {
	conn := testDB(t)

	AccountInsert(conn, "savings", "bank")
	acct, err := AccountGetByName(conn, "savings")
	if err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if acct.Name != "savings" {
		t.Errorf("expected 'savings', got %q", acct.Name)
	}

	_, err = AccountGetByName(conn, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent account")
	}
}

func TestAccountGetFirst(t *testing.T) {
	conn := testDB(t)

	// No accounts yet
	_, err := AccountGetFirst(conn)
	if err == nil {
		t.Error("expected error when no accounts")
	}

	AccountInsert(conn, "first", "bank")
	AccountInsert(conn, "second", "cash")

	acct, err := AccountGetFirst(conn)
	if err != nil {
		t.Fatalf("get first: %v", err)
	}
	if acct.Name != "first" {
		t.Errorf("expected 'first', got %q", acct.Name)
	}
}

func TestAccountEdit(t *testing.T) {
	conn := testDB(t)

	acct, _ := AccountInsert(conn, "old", "bank")
	edited, err := AccountEdit(conn, acct.ID, "new", "wallet")
	if err != nil {
		t.Fatalf("edit: %v", err)
	}
	if edited.Name != "new" || edited.Type != "wallet" {
		t.Errorf("expected new/wallet, got %s/%s", edited.Name, edited.Type)
	}
}

func TestAccountInvalidType(t *testing.T) {
	conn := testDB(t)

	_, err := AccountInsert(conn, "bad", "invalid_type")
	if err == nil {
		t.Error("expected error for invalid account type")
	}
}
