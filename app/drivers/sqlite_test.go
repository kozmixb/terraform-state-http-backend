package drivers

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestSqLiteDriverStoresQuotedPayload(t *testing.T) {
	skipWithoutCGO(t)
	t.Chdir(t.TempDir())

	driver := NewSqLiteDriver()
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
	skipWithoutCGO(t)
	t.Chdir(t.TempDir())

	driver := NewSqLiteDriver()
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

func skipWithoutCGO(t *testing.T) {
	t.Helper()

	output, err := exec.Command("go", "env", "CGO_ENABLED").Output()
	if err != nil {
		t.Skipf("could not determine CGO support: %v", err)
	}
	if strings.TrimSpace(string(output)) == "0" {
		t.Skip("go-sqlite3 requires CGO")
	}
}
