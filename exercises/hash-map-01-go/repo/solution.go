package main

// MyHashMap is a hash map for non-negative integer keys and values,
// built without Go's own map type.
type MyHashMap struct {
}

func NewMyHashMap() *MyHashMap {
	return &MyHashMap{}
}

// Put inserts key with value, or updates it if key already exists.
func (m *MyHashMap) Put(key, value int) {
	// TODO: implement
}

// Get returns the value for key, or -1 if key is absent.
func (m *MyHashMap) Get(key int) int {
	// TODO: implement
	return -1
}

// Remove deletes key if present; does nothing otherwise.
func (m *MyHashMap) Remove(key int) {
	// TODO: implement
}
