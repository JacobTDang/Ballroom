package main

import "fmt"

// BeginState distinguishes what the caller should do next.
type BeginState int

const (
	BeginExecute BeginState = iota
	BeginInFlight
	BeginReplay
)

// BeginResult is what Begin hands back. Response is only meaningful
// when State == BeginReplay.
type BeginResult struct {
	State    BeginState
	Response string
}

type record struct {
	fingerprint string
	inFlight    bool
	deadline    int64 // now_ms past this: treat the record as gone
	response    string
}

// IdempotencyStore: one record per key -- fingerprint, in-flight or
// completed, a deadline that Begin sets and Complete renews, and
// (once completed) the stored response. Past its deadline, a record
// is treated as if it never existed: Begin starts clean, and Complete
// has nothing to attach to.
type IdempotencyStore struct {
	ttlMillis int64
	records   map[string]*record
}

func NewIdempotencyStore(ttlMillis int64) *IdempotencyStore {
	return &IdempotencyStore{ttlMillis: ttlMillis, records: make(map[string]*record)}
}

func (s *IdempotencyStore) BeginAt(key, fingerprint string, nowMillis int64) (BeginResult, error) {
	r, ok := s.records[key]
	if !ok || nowMillis >= r.deadline {
		s.records[key] = &record{
			fingerprint: fingerprint,
			inFlight:    true,
			deadline:    nowMillis + s.ttlMillis,
		}
		return BeginResult{State: BeginExecute}, nil
	}

	if r.fingerprint != fingerprint {
		return BeginResult{}, fmt.Errorf("idempotency: fingerprint conflict for key %q", key)
	}

	if r.inFlight {
		return BeginResult{State: BeginInFlight}, nil
	}
	return BeginResult{State: BeginReplay, Response: r.response}, nil
}

func (s *IdempotencyStore) CompleteAt(key, response string, nowMillis int64) error {
	r, ok := s.records[key]
	if !ok || nowMillis >= r.deadline || !r.inFlight {
		return fmt.Errorf("idempotency: no in-flight request for key %q", key)
	}
	r.inFlight = false
	r.response = response
	r.deadline = nowMillis + s.ttlMillis
	return nil
}
