package main

import "testing"

func makeRecords(ids []int64) []Record {
	out := make([]Record, len(ids))
	for i, id := range ids {
		out[i] = Record{ID: id, Name: "r"}
	}
	return out
}

func tamper(token string) string {
	mid := len(token) / 2
	b := []byte(token)
	if b[mid] != 'X' {
		b[mid] = 'X'
	} else {
		b[mid] = 'Y'
	}
	return string(b)
}

func TestExactlyOnceFullWalk(t *testing.T) {
	ids := make([]int64, 23)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	store := NewCursorStore(makeRecords(ids))

	seen := map[int64]int{}
	var order []int64
	token := ""
	for {
		page, next, err := store.List(5, token)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		for _, r := range page {
			seen[r.ID]++
			order = append(order, r.ID)
		}
		token = next
		if token == "" {
			break
		}
	}
	if len(order) != len(ids) {
		t.Fatalf("walk returned %d records, want %d", len(order), len(ids))
	}
	for _, id := range ids {
		if seen[id] != 1 {
			t.Fatalf("id %d seen %d times, want exactly 1", id, seen[id])
		}
	}
}

func TestEmptyFinalTokenAndNonEmptyWhenMoreRemain(t *testing.T) {
	store := NewCursorStore(makeRecords([]int64{1, 2, 3}))
	page, token, err := store.List(10, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(page) != 3 {
		t.Fatalf("page len = %d, want 3", len(page))
	}
	if token != "" {
		t.Fatal("nothing left, but next_page_token wasn't empty")
	}

	store2 := NewCursorStore(makeRecords([]int64{1, 2, 3, 4, 5}))
	page2, token2, err := store2.List(2, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(page2) != 2 {
		t.Fatalf("page len = %d, want 2", len(page2))
	}
	if token2 == "" {
		t.Fatal("more records remain, but next_page_token was empty")
	}
}

func TestTamperedTokenErrors(t *testing.T) {
	ids := make([]int64, 9)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	store := NewCursorStore(makeRecords(ids))
	_, token, err := store.List(3, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if _, _, err := store.List(3, tamper(token)); err == nil {
		t.Fatal("a tampered page_token was silently accepted")
	}
}

func TestPageSizeClamps(t *testing.T) {
	ids := make([]int64, 60)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	store := NewCursorStore(makeRecords(ids))

	page, _, err := store.List(0, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(page) != 10 {
		t.Fatalf("page_size<=0: got %d records, want the 10-record default", len(page))
	}

	page, _, err = store.List(-5, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(page) != 10 {
		t.Fatalf("negative page_size: got %d records, want the 10-record default", len(page))
	}

	page, _, err = store.List(10000, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(page) != 50 {
		t.Fatalf("oversized page_size: got %d records, want the 50-record max", len(page))
	}
}

func TestParamChangeInvalidatesToken(t *testing.T) {
	ids := make([]int64, 29)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	store := NewCursorStore(makeRecords(ids))
	_, token, err := store.List(5, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if _, _, err := store.List(7, token); err == nil {
		t.Fatal("resuming with a different page_size was silently honored")
	}
}

func TestInsertMidWalkNeverDuplicatesOrSkips(t *testing.T) {
	seedIDs := make([]int64, 20)
	for i := range seedIDs {
		seedIDs[i] = int64((i + 1) * 10) // 10, 20, ..., 200
	}
	store := NewCursorStore(makeRecords(seedIDs))

	page1, token, err := store.List(5, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(page1) != 5 || page1[0].ID != 10 || page1[4].ID != 50 {
		t.Fatalf("unexpected first page: %+v", page1)
	}

	// Lands before everything returned so far -- the classic case where
	// offset pagination re-shows an already-seen record.
	if err := store.Insert(Record{ID: 5, Name: "new-before-cursor"}); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	// Lands past the end of the original seed.
	if err := store.Insert(Record{ID: 999, Name: "new-after-cursor"}); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	seen := map[int64]int{}
	for _, r := range page1 {
		seen[r.ID]++
	}
	for {
		page, next, err := store.List(5, token)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		for _, r := range page {
			seen[r.ID]++
		}
		token = next
		if token == "" {
			break
		}
	}

	for _, id := range seedIDs {
		if seen[id] != 1 {
			t.Fatalf("original id %d seen %d times, want exactly 1", id, seen[id])
		}
	}
}
