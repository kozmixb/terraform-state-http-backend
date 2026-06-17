package drivers

import "testing"

func TestSameLockIDUsesTerraformID(t *testing.T) {
	current := []byte(`{"ID":"same","Operation":"apply"}`)
	requested := []byte(`{"ID":"same","Operation":"force-unlock"}`)

	if !sameLockID(current, requested) {
		t.Fatal("expected matching lock IDs")
	}
}

func TestSameLockIDRejectsDifferentTerraformID(t *testing.T) {
	current := []byte(`{"ID":"current"}`)
	requested := []byte(`{"ID":"requested"}`)

	if sameLockID(current, requested) {
		t.Fatal("expected different lock IDs not to match")
	}
}

func TestSameLockIDFallsBackToTrimmedPayloadComparison(t *testing.T) {
	current := []byte(" raw-lock \n")
	requested := []byte("raw-lock")

	if !sameLockID(current, requested) {
		t.Fatal("expected equivalent raw lock payloads to match")
	}
}

func TestLockIDReturnsEmptyForInvalidPayload(t *testing.T) {
	if got := lockID([]byte("not-json")); got != "" {
		t.Fatalf("expected empty lock ID, got %q", got)
	}
}
