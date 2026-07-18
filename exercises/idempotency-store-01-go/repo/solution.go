package main

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

// IdempotencyStore: tracks one request per key through its lifecycle.
//
// TODO: no fingerprint tracking, no in-flight/completed distinction,
// no deadline at all -- everything after the first Begin just replays
// whatever was last stored, even if nothing ever completed. Every
// rule in the problem statement is still yours to build.
type IdempotencyStore struct {
	ttlMillis int64
	seen      map[string]string // key -> response ("" if never completed)
}

func NewIdempotencyStore(ttlMillis int64) *IdempotencyStore {
	return &IdempotencyStore{ttlMillis: ttlMillis, seen: make(map[string]string)}
}

func (s *IdempotencyStore) BeginAt(key, fingerprint string, nowMillis int64) (BeginResult, error) {
	if _, ok := s.seen[key]; !ok {
		s.seen[key] = ""
		return BeginResult{State: BeginExecute}, nil
	}
	return BeginResult{State: BeginReplay, Response: s.seen[key]}, nil
}

func (s *IdempotencyStore) CompleteAt(key, response string, nowMillis int64) error {
	s.seen[key] = response
	return nil
}
