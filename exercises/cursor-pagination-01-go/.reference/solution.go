package main

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Record is a single row in the store, ordered by ID.
type Record struct {
	ID   int64
	Name string
}

const (
	defaultPageSize = 10
	maxPageSize     = 50
)

func clampPageSize(pageSize int) int {
	switch {
	case pageSize <= 0:
		return defaultPageSize
	case pageSize > maxPageSize:
		return maxPageSize
	default:
		return pageSize
	}
}

func checksum(payload string) string {
	var h int64
	for _, ch := range payload {
		h = (h*131 + int64(ch)) % 1_000_000_007
	}
	return strconv.FormatInt(h, 16)
}

func encodeToken(lastID int64, pageSize int) string {
	payload := fmt.Sprintf("%d:%d", lastID, pageSize)
	raw := payload + ":" + checksum(payload)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func decodeToken(token string) (lastID int64, pageSize int, err error) {
	rawBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid page token %q: %w", token, err)
	}
	parts := strings.Split(string(rawBytes), ":")
	if len(parts) != 3 {
		return 0, 0, fmt.Errorf("invalid page token %q", token)
	}
	payload := parts[0] + ":" + parts[1]
	if checksum(payload) != parts[2] {
		return 0, 0, fmt.Errorf("invalid page token %q: checksum mismatch", token)
	}
	lastID, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid page token %q", token)
	}
	pageSize, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid page token %q", token)
	}
	return lastID, pageSize, nil
}

// CursorStore: keyset pagination. The token names the last id seen
// (plus the page_size it was issued for), never a raw offset -- so a
// walk can't be thrown off by inserts, and resuming with a different
// page size is a loud error rather than a silent behavior change.
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
	if _, exists := s.records[r.ID]; exists {
		return fmt.Errorf("duplicate id: %d", r.ID)
	}
	s.records[r.ID] = r
	return nil
}

func (s *CursorStore) List(pageSize int, pageToken string) ([]Record, string, error) {
	effective := clampPageSize(pageSize)

	var cursorID int64
	hasCursor := false
	if pageToken != "" {
		id, encodedSize, err := decodeToken(pageToken)
		if err != nil {
			return nil, "", err
		}
		if encodedSize != effective {
			return nil, "", fmt.Errorf("page_token was issued for page_size=%d, not %d", encodedSize, effective)
		}
		cursorID = id
		hasCursor = true
	}

	candidates := make([]Record, 0, len(s.records))
	for _, r := range s.records {
		if !hasCursor || r.ID > cursorID {
			candidates = append(candidates, r)
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].ID < candidates[j].ID })

	pageLen := effective
	if pageLen > len(candidates) {
		pageLen = len(candidates)
	}
	page := candidates[:pageLen]

	nextToken := ""
	if len(candidates) > pageLen {
		nextToken = encodeToken(page[len(page)-1].ID, effective)
	}
	return page, nextToken, nil
}
