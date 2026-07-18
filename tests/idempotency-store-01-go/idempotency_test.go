package main

import "testing"

const ttl = int64(1000)

func TestLifecycleMatrix(t *testing.T) {
	store := NewIdempotencyStore(ttl)

	res, err := store.BeginAt("k1", "fp-a", 0)
	if err != nil || res.State != BeginExecute {
		t.Fatalf("first Begin: got %+v, %v, want BeginExecute, nil", res, err)
	}

	// A second Begin while the first is still running (same fingerprint)
	// must not re-execute -- it's a duplicate in flight.
	res, err = store.BeginAt("k1", "fp-a", 10)
	if err != nil || res.State != BeginInFlight {
		t.Fatalf("duplicate Begin: got %+v, %v, want BeginInFlight, nil", res, err)
	}

	if err := store.CompleteAt("k1", "RESULT-1", 20); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	res, err = store.BeginAt("k1", "fp-a", 30)
	if err != nil || res.State != BeginReplay || res.Response != "RESULT-1" {
		t.Fatalf("replay Begin: got %+v, %v, want BeginReplay/RESULT-1, nil", res, err)
	}
}

func TestByteIdenticalReplay(t *testing.T) {
	store := NewIdempotencyStore(ttl)
	store.BeginAt("k1", "fp-a", 0)
	payload := `{"amount": 4200, "currency": "usd", "note": "café"}`
	if err := store.CompleteAt("k1", payload, 10); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	res, err := store.BeginAt("k1", "fp-a", 20)
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if res.Response != payload {
		t.Fatalf("replay must return the stored response byte-for-byte: got %q, want %q", res.Response, payload)
	}
}

func TestConflictOnLiveKeyWithDifferentFingerprint(t *testing.T) {
	store := NewIdempotencyStore(ttl)
	store.BeginAt("k1", "fp-a", 0)
	if _, err := store.BeginAt("k1", "fp-b", 5); err == nil {
		t.Fatal("in-flight key with a mismatched fingerprint didn't error")
	}

	store2 := NewIdempotencyStore(ttl)
	store2.BeginAt("k2", "fp-a", 0)
	store2.CompleteAt("k2", "RESULT", 5)
	if _, err := store2.BeginAt("k2", "fp-b", 10); err == nil {
		t.Fatal("completed key with a mismatched fingerprint didn't error")
	}
}

func TestExactTTLBoundary(t *testing.T) {
	store := NewIdempotencyStore(ttl)
	store.BeginAt("k1", "fp-a", 0)
	store.CompleteAt("k1", "RESULT", 0) // deadline = 0 + ttl

	res, err := store.BeginAt("k1", "fp-a", ttl-1)
	if err != nil || res.State != BeginReplay || res.Response != "RESULT" {
		t.Fatalf("just inside the retention window: got %+v, %v, want BeginReplay/RESULT", res, err)
	}

	res, err = store.BeginAt("k1", "fp-a", ttl)
	if err != nil || res.State != BeginExecute {
		t.Fatalf("deadline reached: got %+v, %v, want BeginExecute (brand new key)", res, err)
	}
}

func TestCompleteOnUnknownOrExpiredErrors(t *testing.T) {
	store := NewIdempotencyStore(ttl)
	if err := store.CompleteAt("never-begun", "R", 5); err == nil {
		t.Fatal("Complete on an unknown key didn't error")
	}

	store2 := NewIdempotencyStore(ttl)
	store2.BeginAt("k1", "fp-a", 0) // in-flight, deadline = ttl
	if err := store2.CompleteAt("k1", "R", ttl); err == nil {
		t.Fatal("Complete on a key whose deadline already passed didn't error")
	}
}
