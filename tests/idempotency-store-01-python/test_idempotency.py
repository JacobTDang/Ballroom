import pytest

from solution import IdempotencyStore

TTL = 1000


def test_lifecycle_matrix():
    store = IdempotencyStore(TTL)

    state, resp = store.begin_at("k1", "fp-a", 0)
    assert state == "execute" and resp is None

    # A second Begin while the first is still running (same fingerprint)
    # must not re-execute -- it's a duplicate in flight.
    state, resp = store.begin_at("k1", "fp-a", 10)
    assert state == "in-flight"

    store.complete_at("k1", "RESULT-1", 20)

    # A Begin with the same key/fingerprint now replays the stored result.
    state, resp = store.begin_at("k1", "fp-a", 30)
    assert state == "replay" and resp == "RESULT-1"


def test_byte_identical_replay():
    store = IdempotencyStore(TTL)
    store.begin_at("k1", "fp-a", 0)
    payload = '{"amount": 4200, "currency": "usd", "note": "café"}'
    store.complete_at("k1", payload, 10)

    _, resp = store.begin_at("k1", "fp-a", 20)
    assert resp == payload, "replay must return the stored response byte-for-byte"


def test_conflict_on_live_key_with_different_fingerprint():
    # Conflict while in-flight.
    store = IdempotencyStore(TTL)
    store.begin_at("k1", "fp-a", 0)
    with pytest.raises(ValueError):
        store.begin_at("k1", "fp-b", 5)

    # Conflict on a completed (but not expired) key.
    store2 = IdempotencyStore(TTL)
    store2.begin_at("k2", "fp-a", 0)
    store2.complete_at("k2", "RESULT", 5)
    with pytest.raises(ValueError):
        store2.begin_at("k2", "fp-b", 10)


def test_exact_ttl_boundary():
    store = IdempotencyStore(TTL)
    store.begin_at("k1", "fp-a", 0)
    store.complete_at("k1", "RESULT", 0)  # deadline = 0 + TTL

    state, resp = store.begin_at("k1", "fp-a", TTL - 1)
    assert state == "replay" and resp == "RESULT", "still inside the retention window"

    state, resp = store.begin_at("k1", "fp-a", TTL)
    assert state == "execute", "deadline reached -- must be treated as a brand new key"


def test_complete_on_unknown_or_expired_errors():
    store = IdempotencyStore(TTL)
    with pytest.raises(ValueError):
        store.complete_at("never-begun", "R", 5)

    store2 = IdempotencyStore(TTL)
    store2.begin_at("k1", "fp-a", 0)  # in-flight, deadline = TTL
    with pytest.raises(ValueError):
        store2.complete_at("k1", "R", TTL)  # deadline reached before it ever completed
