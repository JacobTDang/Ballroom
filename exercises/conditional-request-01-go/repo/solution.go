package main

// GetResult is what Get hands back. ETag/Body are only meaningful
// when Status == 200.
type GetResult struct {
	Status int
	ETag   string
	Body   string
}

// PutResult is what Put hands back. ETag is only meaningful when
// Status == 200.
type PutResult struct {
	Status int
	ETag   string
}

// ConditionalStore: a versioned key-value store.
//
// TODO: no versioning at all -- Put always succeeds and overwrites
// unconditionally (the classic lost-update bug this store exists to
// prevent), and Get hands out the same constant etag forever. Every
// rule in the problem statement is still yours to build.
type ConditionalStore struct {
	items map[string]string
}

func NewConditionalStore() *ConditionalStore {
	return &ConditionalStore{items: make(map[string]string)}
}

func (s *ConditionalStore) Get(key, ifNoneMatch string) GetResult {
	body, ok := s.items[key]
	if !ok {
		return GetResult{Status: 404}
	}
	return GetResult{Status: 200, ETag: "1", Body: body}
}

func (s *ConditionalStore) Put(key, body, ifMatch string) PutResult {
	s.items[key] = body
	return PutResult{Status: 200, ETag: "1"}
}

func (s *ConditionalStore) Delete(key string) {
	delete(s.items, key)
}
