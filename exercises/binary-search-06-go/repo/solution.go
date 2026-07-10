package main

// TimeMap stores multiple values per key, each tagged with the
// timestamp it was set at.
type TimeMap struct {
}

func NewTimeMap() *TimeMap {
	return &TimeMap{}
}

func (m *TimeMap) Set(key, value string, timestamp int) {
	// TODO: implement
}

func (m *TimeMap) Get(key string, timestamp int) string {
	// TODO: implement
	return ""
}
