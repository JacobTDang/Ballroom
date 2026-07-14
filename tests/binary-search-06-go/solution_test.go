package main

import "testing"

func TestTimeMap(t *testing.T) {
	m := NewTimeMap()
	m.Set("foo", "bar", 1)
	if got := m.Get("foo", 1); got != "bar" {
		t.Fatalf("Get(foo, 1) = %q, want bar", got)
	}
	if got := m.Get("foo", 3); got != "bar" {
		t.Fatalf("Get(foo, 3) = %q, want bar (no exact match, falls back to timestamp 1)", got)
	}
	m.Set("foo", "bar2", 4)
	if got := m.Get("foo", 4); got != "bar2" {
		t.Fatalf("Get(foo, 4) = %q, want bar2", got)
	}
	if got := m.Get("foo", 5); got != "bar2" {
		t.Fatalf("Get(foo, 5) = %q, want bar2", got)
	}
}

func TestTimeMap_GetBeforeAnySetReturnsEmpty(t *testing.T) {
	m := NewTimeMap()
	m.Set("foo", "bar", 5)
	if got := m.Get("foo", 1); got != "" {
		t.Fatalf("Get(foo, 1) = %q, want empty (timestamp before first set)", got)
	}
}

func TestTimeMap_UnknownKeyReturnsEmpty(t *testing.T) {
	m := NewTimeMap()
	if got := m.Get("missing", 1); got != "" {
		t.Fatalf("Get(missing, 1) = %q, want empty", got)
	}
}

func TestTimeMap_MultipleKeysDoNotInterfere(t *testing.T) {
	m := NewTimeMap()
	m.Set("foo", "foo-val", 1)
	m.Set("bar", "bar-val", 2)
	if got := m.Get("foo", 10); got != "foo-val" {
		t.Fatalf("Get(foo, 10) = %q, want foo-val", got)
	}
	if got := m.Get("bar", 10); got != "bar-val" {
		t.Fatalf("Get(bar, 10) = %q, want bar-val", got)
	}
}
