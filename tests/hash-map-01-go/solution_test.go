package main

import "testing"

func TestPutThenGetRoundtrip(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1, 100)
	if got := m.Get(1); got != 100 {
		t.Errorf("Get(1) = %d, want 100", got)
	}
}

func TestGetMissingReturnsMinusOne(t *testing.T) {
	m := NewMyHashMap()
	if got := m.Get(42); got != -1 {
		t.Errorf("Get(42) = %d, want -1", got)
	}
}

func TestPutOverwritesExistingKey(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1, 100)
	m.Put(1, 200)
	if got := m.Get(1); got != 200 {
		t.Errorf("Get(1) = %d, want 200 after overwrite", got)
	}
}

func TestRemoveThenGetReturnsMinusOne(t *testing.T) {
	m := NewMyHashMap()
	m.Put(7, 70)
	m.Remove(7)
	if got := m.Get(7); got != -1 {
		t.Errorf("Get(7) = %d, want -1 after Remove", got)
	}
}

func TestRemoveMissingKeyIsANoop(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1, 10)
	m.Remove(99)
	if got := m.Get(1); got != 10 {
		t.Errorf("Get(1) = %d, want 10 -- removing a missing key must not disturb others", got)
	}
}

func TestCollidingKeysAreKeptSeparate(t *testing.T) {
	m := NewMyHashMap()
	keys := []int{1, 1025, 2049, 1001, 2001}
	for _, k := range keys {
		m.Put(k, k*3)
	}
	for _, k := range keys {
		if got := m.Get(k); got != k*3 {
			t.Errorf("Get(%d) = %d, want %d -- colliding keys must chain", k, got, k*3)
		}
	}
}

func TestRemovingOneCollidingKeyKeepsTheOthers(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1, 11)
	m.Put(1025, 22)
	m.Put(2049, 33)
	m.Remove(1025)
	if got := m.Get(1); got != 11 {
		t.Errorf("Get(1) = %d, want 11", got)
	}
	if got := m.Get(1025); got != -1 {
		t.Errorf("Get(1025) = %d, want -1", got)
	}
	if got := m.Get(2049); got != 33 {
		t.Errorf("Get(2049) = %d, want 33", got)
	}
}

func TestZeroKeyAndZeroValue(t *testing.T) {
	m := NewMyHashMap()
	m.Put(0, 0)
	if got := m.Get(0); got != 0 {
		t.Errorf("Get(0) = %d, want 0", got)
	}
}

func TestLargeKeyBounds(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1000000, 999)
	if got := m.Get(1000000); got != 999 {
		t.Errorf("Get(1000000) = %d, want 999", got)
	}
}

func TestManyKeysNoAliasing(t *testing.T) {
	m := NewMyHashMap()
	for k := 0; k < 500; k++ {
		m.Put(k, k*2)
	}
	for k := 0; k < 500; k++ {
		if got := m.Get(k); got != k*2 {
			t.Fatalf("Get(%d) = %d, want %d", k, got, k*2)
		}
	}
}

func TestInterleavedPutRemoveSequence(t *testing.T) {
	m := NewMyHashMap()
	m.Put(5, 1)
	m.Put(6, 2)
	m.Remove(5)
	m.Put(6, 3)
	m.Put(5, 4)
	if got := m.Get(5); got != 4 {
		t.Errorf("Get(5) = %d, want 4", got)
	}
	if got := m.Get(6); got != 3 {
		t.Errorf("Get(6) = %d, want 3", got)
	}
}
