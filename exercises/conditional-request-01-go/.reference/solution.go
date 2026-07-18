package main

import "strconv"

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

type entry struct {
	etag string
	body string
}

// ConditionalStore: every successful write (create or update) draws
// its etag from one store-wide monotonic sequence -- so deleting a
// key and recreating it never reissues an old etag, and a stale
// cached etag can never falsely match again, on either the read or
// write side.
type ConditionalStore struct {
	items map[string]*entry
	seq   int64
}

func NewConditionalStore() *ConditionalStore {
	return &ConditionalStore{items: make(map[string]*entry)}
}

func (s *ConditionalStore) nextETag() string {
	s.seq++
	return strconv.FormatInt(s.seq, 10)
}

func (s *ConditionalStore) Get(key, ifNoneMatch string) GetResult {
	e, ok := s.items[key]
	if !ok {
		return GetResult{Status: 404}
	}
	if ifNoneMatch != "" && ifNoneMatch == e.etag {
		return GetResult{Status: 304}
	}
	return GetResult{Status: 200, ETag: e.etag, Body: e.body}
}

func (s *ConditionalStore) Put(key, body, ifMatch string) PutResult {
	e, ok := s.items[key]
	if !ok {
		if ifMatch != "" {
			// Claimed a version of something that doesn't exist.
			return PutResult{Status: 412}
		}
		etag := s.nextETag()
		s.items[key] = &entry{etag: etag, body: body}
		return PutResult{Status: 200, ETag: etag}
	}

	if ifMatch == "" {
		// Blind overwrite of an existing resource: refused on purpose.
		return PutResult{Status: 428}
	}
	if ifMatch != e.etag {
		// Stale precondition -- state is left exactly as it was.
		return PutResult{Status: 412}
	}

	etag := s.nextETag()
	e.etag = etag
	e.body = body
	return PutResult{Status: 200, ETag: etag}
}

func (s *ConditionalStore) Delete(key string) {
	delete(s.items, key)
}
