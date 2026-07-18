from solution import ConditionalStore


def test_get_missing_returns_404():
    store = ConditionalStore()
    status, etag, body = store.get("missing")
    assert status == 404
    assert etag is None and body is None


def test_put_create_flows():
    store = ConditionalStore()

    status, etag = store.put("k1", "v1")
    assert status == 200
    assert etag

    store2 = ConditionalStore()
    status, etag = store2.put("k1", "v1", if_match="some-version")
    assert status == 412, "claimed a version of a resource that doesn't exist yet"


def test_put_update_requires_if_match():
    store = ConditionalStore()
    _, etag1 = store.put("k1", "v1")

    status, etag = store.put("k1", "v2")  # no if_match on an existing key
    assert status == 428
    assert etag is None

    # State must be untouched.
    get_status, get_etag, get_body = store.get("k1")
    assert (get_status, get_etag, get_body) == (200, etag1, "v1")


def test_put_update_stale_if_match_412_leaves_state_unchanged():
    store = ConditionalStore()
    _, etag1 = store.put("k1", "v1")

    status, etag = store.put("k1", "v2", if_match="not-" + etag1)
    assert status == 412
    assert etag is None

    # A failed conditional write must never partially apply.
    get_status, get_etag, get_body = store.get("k1")
    assert (get_status, get_etag, get_body) == (200, etag1, "v1")


def test_put_update_correct_if_match_succeeds_and_rotates_etag():
    store = ConditionalStore()
    _, etag1 = store.put("k1", "v1")

    status, etag2 = store.put("k1", "v2", if_match=etag1)
    assert status == 200
    assert etag2 and etag2 != etag1

    get_status, get_etag, get_body = store.get("k1")
    assert (get_status, get_etag, get_body) == (200, etag2, "v2")


def test_get_if_none_match_matrix():
    store = ConditionalStore()
    _, etag = store.put("k1", "v1")

    status, _, _ = store.get("k1", if_none_match=etag)
    assert status == 304, "matching if_none_match must 304"

    status, got_etag, got_body = store.get("k1", if_none_match="stale-etag")
    assert status == 200 and got_etag == etag and got_body == "v1"


def test_no_etag_resurrection_after_recreate():
    store = ConditionalStore()
    _, etag1 = store.put("b", "first")
    store.delete("b")
    _, etag2 = store.put("b", "first")  # same body, recreated key

    assert etag2 != etag1, "a recreated resource reused its old etag"

    # Read side: the old etag must not falsely 304 the new resource.
    status, _, _ = store.get("b", if_none_match=etag1)
    assert status == 200, "a stale pre-delete etag falsely matched the recreated resource"

    # Write side: the old etag must not succeed as a precondition either.
    status, _ = store.put("b", "second", if_match=etag1)
    assert status == 412, "a stale pre-delete etag falsely satisfied If-Match"
