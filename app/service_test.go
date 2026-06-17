package app

import (
	"errors"
	"terraform-state-http-backend/app/drivers"
	"testing"
)

func TestGetDefaultUsesFileDriverWhenUnset(t *testing.T) {
	t.Setenv("DRIVER", "")

	if got := getDefault(); got != "file" {
		t.Fatalf("expected default driver file, got %q", got)
	}
}

func TestGetDefaultUsesConfiguredDriver(t *testing.T) {
	t.Setenv("DRIVER", "sqlite")

	if got := getDefault(); got != "sqlite" {
		t.Fatalf("expected configured driver sqlite, got %q", got)
	}
}

func TestServiceUpdateShowAndMissingState(t *testing.T) {
	resetDriver(t)
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	payload := []byte(`{"version":4}`)

	result, err := Update("group", "key", payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(payload) {
		t.Fatalf("expected update result %q, got %q", payload, result)
	}

	result, err = Show("group", "key")
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(payload) {
		t.Fatalf("expected show result %q, got %q", payload, result)
	}

	_, err = Show("group", "missing")
	if !errors.Is(err, drivers.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceLockAndUnlock(t *testing.T) {
	resetDriver(t)
	t.Setenv("DRIVER", "file")
	t.Chdir(t.TempDir())

	firstLock := []byte(`{"ID":"first"}`)
	secondLock := []byte(`{"ID":"second"}`)

	if _, err := Lock("group", "key", firstLock); err != nil {
		t.Fatal(err)
	}

	current, err := Lock("group", "key", secondLock)
	if !errors.Is(err, drivers.ErrLocked) {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
	if string(current) != string(firstLock) {
		t.Fatalf("expected current lock %q, got %q", firstLock, current)
	}

	if err := Unlock("group", "key", secondLock); !errors.Is(err, drivers.ErrUnlockMismatch) {
		t.Fatalf("expected ErrUnlockMismatch, got %v", err)
	}
	if err := Unlock("group", "key", firstLock); err != nil {
		t.Fatal(err)
	}
	if _, err := Lock("group", "key", secondLock); err != nil {
		t.Fatal(err)
	}
}

func resetDriver(t *testing.T) {
	t.Helper()

	driverLock.Lock()
	driver = nil
	driverLock.Unlock()
}
