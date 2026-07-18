package main

import "testing"

func kv(records []*Record) []Record {
	out := make([]Record, len(records))
	for i, r := range records {
		out[i] = *r
	}
	return out
}

func recordsEqual(got, want []Record) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestDedupe(t *testing.T) {
	t.Run("value-equal but distinct objects collapse", func(t *testing.T) {
		r1 := &Record{"a", 1}
		r2 := &Record{"a", 1} // distinct pointer, same value
		r3 := &Record{"b", 2}
		got := kv(Dedupe([]*Record{r1, r2, r3}))
		want := []Record{{"a", 1}, {"b", 2}}
		if !recordsEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("triple duplicate collapses to one", func(t *testing.T) {
		r1 := &Record{"x", 5}
		r2 := &Record{"x", 5}
		r3 := &Record{"x", 5}
		got := kv(Dedupe([]*Record{r1, r2, r3}))
		want := []Record{{"x", 5}}
		if !recordsEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("all distinct records survive", func(t *testing.T) {
		r1 := &Record{"a", 1}
		r2 := &Record{"b", 2}
		r3 := &Record{"c", 3}
		got := kv(Dedupe([]*Record{r1, r2, r3}))
		want := []Record{{"a", 1}, {"b", 2}, {"c", 3}}
		if !recordsEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("interleaved duplicates preserve first occurrence order", func(t *testing.T) {
		r1 := &Record{"a", 1}
		r2 := &Record{"b", 2}
		r3 := &Record{"a", 1}
		r4 := &Record{"c", 3}
		r5 := &Record{"b", 2}
		got := kv(Dedupe([]*Record{r1, r2, r3, r4, r5}))
		want := []Record{{"a", 1}, {"b", 2}, {"c", 3}}
		if !recordsEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("same key different value is not a duplicate", func(t *testing.T) {
		r1 := &Record{"a", 1}
		r2 := &Record{"a", 2}
		got := kv(Dedupe([]*Record{r1, r2}))
		want := []Record{{"a", 1}, {"a", 2}}
		if !recordsEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
