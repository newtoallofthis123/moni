package db

import (
	"testing"
)

func TestBucketCreateAndList(t *testing.T) {
	conn := testDB(t)

	b, err := BucketInsert(conn, "vacation", 5000)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if b.Name != "vacation" || b.Target != 5000 || b.Current != 0 {
		t.Errorf("unexpected bucket: %+v", b)
	}

	buckets, _ := BucketList(conn)
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
}

func TestBucketAddFunds(t *testing.T) {
	conn := testDB(t)

	b, _ := BucketInsert(conn, "emergency", 10000)
	updated, err := BucketAddFunds(conn, b.ID, 2500)
	if err != nil {
		t.Fatalf("add funds: %v", err)
	}
	if updated.Current != 2500 {
		t.Errorf("expected 2500, got %.2f", updated.Current)
	}

	// Add more
	updated, _ = BucketAddFunds(conn, b.ID, 1500)
	if updated.Current != 4000 {
		t.Errorf("expected 4000, got %.2f", updated.Current)
	}
}

func TestBucketGetByName(t *testing.T) {
	conn := testDB(t)

	BucketInsert(conn, "car", 20000)
	b, err := BucketGetByName(conn, "car")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if b.Name != "car" {
		t.Errorf("expected 'car', got %q", b.Name)
	}

	_, err = BucketGetByName(conn, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent bucket")
	}
}

func TestBucketDuplicateName(t *testing.T) {
	conn := testDB(t)

	BucketInsert(conn, "save", 1000)
	_, err := BucketInsert(conn, "save", 2000)
	if err == nil {
		t.Error("expected error for duplicate bucket name")
	}
}

func TestBucketDelete(t *testing.T) {
	conn := testDB(t)

	b, _ := BucketInsert(conn, "temp", 500)
	if err := BucketDelete(conn, b.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	buckets, _ := BucketList(conn)
	if len(buckets) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(buckets))
	}

	if err := BucketDelete(conn, 9999); err == nil {
		t.Error("expected error deleting nonexistent bucket")
	}
}

func TestBucketEdit(t *testing.T) {
	conn := testDB(t)

	b, _ := BucketInsert(conn, "old", 1000)
	BucketAddFunds(conn, b.ID, 200)

	edited, err := BucketEdit(conn, b.ID, "new", 2000)
	if err != nil {
		t.Fatalf("edit: %v", err)
	}
	if edited.Name != "new" || edited.Target != 2000 {
		t.Errorf("expected new/2000, got %s/%.0f", edited.Name, edited.Target)
	}
	// Current should be preserved
	if edited.Current != 200 {
		t.Errorf("expected current 200 preserved, got %.2f", edited.Current)
	}
}
