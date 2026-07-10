package main

// timeEntry pairs a timestamp with the value set at that time.
type timeEntry struct {
	timestamp int
	value     string
}

// TimeMap stores multiple values per key, each tagged with the
// timestamp it was set at.
type TimeMap struct {
	store map[string][]timeEntry
}

func NewTimeMap() *TimeMap {
	return &TimeMap{store: make(map[string][]timeEntry)}
}

// Set appends value for key — timestamps arrive strictly increasing,
// so the per-key slice stays sorted without needing to insert.
func (m *TimeMap) Set(key, value string, timestamp int) {
	m.store[key] = append(m.store[key], timeEntry{timestamp, value})
}

// Get binary-searches for the entry with the largest timestamp <=
// the query timestamp.
func (m *TimeMap) Get(key string, timestamp int) string {
	entries := m.store[key]
	lo, hi := 0, len(entries)-1
	res := ""
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if entries[mid].timestamp <= timestamp {
			res = entries[mid].value
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return res
}
