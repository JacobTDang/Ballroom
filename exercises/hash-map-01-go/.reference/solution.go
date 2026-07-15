package main

// MyHashMap is a hash map for non-negative integer keys and values,
// built without Go's own map type.
type MyHashMap struct {
	buckets [][]entry
}

type entry struct {
	key, value int
}

const bucketCount = 1024

func NewMyHashMap() *MyHashMap {
	return &MyHashMap{buckets: make([][]entry, bucketCount)}
}

// Put inserts key with value, or updates it if key already exists.
func (m *MyHashMap) Put(key, value int) {
	idx := key % bucketCount
	for i, e := range m.buckets[idx] {
		if e.key == key {
			m.buckets[idx][i].value = value
			return
		}
	}
	m.buckets[idx] = append(m.buckets[idx], entry{key, value})
}

// Get returns the value for key, or -1 if key is absent.
func (m *MyHashMap) Get(key int) int {
	for _, e := range m.buckets[key%bucketCount] {
		if e.key == key {
			return e.value
		}
	}
	return -1
}

// Remove deletes key if present; does nothing otherwise.
func (m *MyHashMap) Remove(key int) {
	idx := key % bucketCount
	for i, e := range m.buckets[idx] {
		if e.key == key {
			m.buckets[idx] = append(m.buckets[idx][:i], m.buckets[idx][i+1:]...)
			return
		}
	}
}
