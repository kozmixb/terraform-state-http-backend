package drivers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileDriverExistsShowAndUpdate(t *testing.T) {
	t.Chdir(t.TempDir())

	driver := FileDriver{}
	exists, err := driver.Exists("group", "key")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("expected state not to exist")
	}

	payload := []byte(`{"version":4}`)
	result, err := driver.Update("group", "key", payload)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(payload) {
		t.Fatalf("expected update result %q, got %q", payload, result)
	}

	exists, err = driver.Exists("group", "key")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected state to exist")
	}

	result, err = driver.Show("group", "key")
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(payload) {
		t.Fatalf("expected show result %q, got %q", payload, result)
	}
}

func TestFileDriverShowMissingState(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := FileDriver{}.Show("group", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFileDriverUnlockMissingLockIsNoop(t *testing.T) {
	t.Chdir(t.TempDir())

	err := FileDriver{}.Unlock("group", "key", []byte(`{"ID":"missing"}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileDriverUsesExpectedStoragePaths(t *testing.T) {
	t.Chdir(t.TempDir())

	payload := []byte(`{"version":4}`)
	if _, err := (FileDriver{}).Update("group", "key", payload); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join("storage", "group-key.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(payload) {
		t.Fatalf("expected file content %q, got %q", payload, content)
	}
}
