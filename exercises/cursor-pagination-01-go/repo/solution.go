package main

import (
	"sort"
	"strconv"
)

// Record is a single row in the store, ordered by ID.
type Record struct {
	ID   int64
	Name string
}

// CursorStore: keyset pagination over id-sorted records.
//
// TODO: this paginates by raw position offset -- it looks fine until
// the data changes between calls. No tamper check, no page-size
// clamp, no page-size binding in the token. Every rule in the problem
// statement is still yours to build.
type CursorStore struct {
	records map[int64]Record
}

func NewCursorStore(records []Record) *CursorStore {
	m := make(map[int64]Record, len(records))
	for _, r := range records {
		m[r.ID] = r
	}
	return &CursorStore{records: m}
}

func (s *CursorStore) Insert(r Record) error {
	s.records[r.ID] = r
	return nil
}

func (s *CursorStore) List(pageSize int, pageToken string) ([]Record, string, error) {
	offset := 0
	if pageToken != "" {
		o, err := strconv.Atoi(pageToken)
		if err != nil {
			return nil, "", err
		}
		offset = o
	}

	ordered := make([]Record, 0, len(s.records))
	for _, r := range s.records {
		ordered = append(ordered, r)
	}
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })

	end := offset + pageSize
	if end > len(ordered) {
		end = len(ordered)
	}
	if offset > len(ordered) {
		offset = len(ordered)
	}
	page := ordered[offset:end]

	nextToken := ""
	if end < len(ordered) {
		nextToken = strconv.Itoa(end)
	}
	return page, nextToken, nil
}
