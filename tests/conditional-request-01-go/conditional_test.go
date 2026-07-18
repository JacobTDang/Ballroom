package main

import "testing"

func TestGetMissingReturns404(t *testing.T) {
	store := NewConditionalStore()
	res := store.Get("missing", "")
	if res.Status != 404 {
		t.Fatalf("Status = %d, want 404", res.Status)
	}
}

func TestPutCreateFlows(t *testing.T) {
	store := NewConditionalStore()
	res := store.Put("k1", "v1", "")
	if res.Status != 200 || res.ETag == "" {
		t.Fatalf("fresh create: got %+v, want 200 with a non-empty etag", res)
	}

	store2 := NewConditionalStore()
	res2 := store2.Put("k1", "v1", "some-version")
	if res2.Status != 412 {
		t.Fatalf("create with if_match set: got %+v, want 412", res2)
	}
}

func TestPutUpdateRequiresIfMatch(t *testing.T) {
	store := NewConditionalStore()
	created := store.Put("k1", "v1", "")

	res := store.Put("k1", "v2", "") // no if_match on an existing key
	if res.Status != 428 {
		t.Fatalf("update without if_match: got %+v, want 428", res)
	}

	get := store.Get("k1", "")
	if get.Status != 200 || get.ETag != created.ETag || get.Body != "v1" {
		t.Fatalf("state changed after a 428: got %+v", get)
	}
}

func TestPutUpdateStaleIfMatch412LeavesStateUnchanged(t *testing.T) {
	store := NewConditionalStore()
	created := store.Put("k1", "v1", "")

	res := store.Put("k1", "v2", "not-"+created.ETag)
	if res.Status != 412 {
		t.Fatalf("stale if_match: got %+v, want 412", res)
	}

	get := store.Get("k1", "")
	if get.Status != 200 || get.ETag != created.ETag || get.Body != "v1" {
		t.Fatalf("a failed conditional write changed state: got %+v", get)
	}
}

func TestPutUpdateCorrectIfMatchSucceedsAndRotatesETag(t *testing.T) {
	store := NewConditionalStore()
	created := store.Put("k1", "v1", "")

	res := store.Put("k1", "v2", created.ETag)
	if res.Status != 200 || res.ETag == "" || res.ETag == created.ETag {
		t.Fatalf("update with correct if_match: got %+v, want 200 with a fresh etag", res)
	}

	get := store.Get("k1", "")
	if get.Status != 200 || get.ETag != res.ETag || get.Body != "v2" {
		t.Fatalf("update didn't take effect: got %+v", get)
	}
}

func TestGetIfNoneMatchMatrix(t *testing.T) {
	store := NewConditionalStore()
	created := store.Put("k1", "v1", "")

	res := store.Get("k1", created.ETag)
	if res.Status != 304 {
		t.Fatalf("matching if_none_match: got %+v, want 304", res)
	}

	res2 := store.Get("k1", "stale-etag")
	if res2.Status != 200 || res2.ETag != created.ETag || res2.Body != "v1" {
		t.Fatalf("stale if_none_match: got %+v, want a fresh 200", res2)
	}
}

func TestNoETagResurrectionAfterRecreate(t *testing.T) {
	store := NewConditionalStore()
	first := store.Put("b", "first", "")
	store.Delete("b")
	second := store.Put("b", "first", "") // same body, recreated key

	if second.ETag == first.ETag {
		t.Fatal("a recreated resource reused its old etag")
	}

	get := store.Get("b", first.ETag)
	if get.Status != 200 {
		t.Fatalf("a stale pre-delete etag falsely matched the recreated resource: got %+v", get)
	}

	put := store.Put("b", "second", first.ETag)
	if put.Status != 412 {
		t.Fatalf("a stale pre-delete etag falsely satisfied If-Match: got %+v", put)
	}
}
