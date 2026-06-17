package drivers

import (
	"errors"
	"testing"
)

func TestSqLiteDriverStoresQuotedPayload(t *testing.T) {
	t.Chdir(t.TempDir())

	driver := NewSqLiteDriver()
	defer driver.db.Close()
	payload := []byte(`{"outputs":{"quote":"it's ok"}}`)

	result, err := driver.Update("group", "key", payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(payload) {
		t.Fatalf("expected update result %q, got %q", payload, result)
	}

	result, err = driver.Show("group", "key")
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(payload) {
		t.Fatalf("expected show result %q, got %q", payload, result)
	}
}

func TestSqLiteDriverLockConflict(t *testing.T) {
	t.Chdir(t.TempDir())

	driver := NewSqLiteDriver()
	defer driver.db.Close()
	firstLock := []byte(`{"ID":"first"}`)
	secondLock := []byte(`{"ID":"second"}`)

	if _, err := driver.Lock("group", "key", firstLock); err != nil {
		t.Fatal(err)
	}

	current, err := driver.Lock("group", "key", secondLock)
	if !errors.Is(err, ErrLocked) {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
	if string(current) != string(firstLock) {
		t.Fatalf("expected current lock %q, got %q", firstLock, current)
	}

	if err := driver.Unlock("group", "key", secondLock); !errors.Is(err, ErrUnlockMismatch) {
		t.Fatalf("expected ErrUnlockMismatch, got %v", err)
	}
	if err := driver.Unlock("group", "key", firstLock); err != nil {
		t.Fatal(err)
	}
	if _, err := driver.Lock("group", "key", secondLock); err != nil {
		t.Fatal(err)
	}
}
