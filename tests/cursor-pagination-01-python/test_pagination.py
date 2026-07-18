import pytest

from solution import CursorStore


def _records(ids):
    return [{"id": i, "name": f"r{i}"} for i in ids]


def _tamper(token):
    mid = len(token) // 2
    ch = "X" if token[mid] != "X" else "Y"
    return token[:mid] + ch + token[mid + 1:]


def test_exactly_once_full_walk():
    ids = list(range(1, 24))  # 23 records
    store = CursorStore(_records(ids))
    seen = []
    token = ""
    while True:
        page, token = store.list(5, token)
        seen.extend(r["id"] for r in page)
        if not token:
            break
    assert sorted(seen) == ids, "walk didn't visit every record exactly once"
    assert len(seen) == len(set(seen)), "walk produced a duplicate"


def test_empty_final_token_and_nonempty_when_more_remain():
    store = CursorStore(_records([1, 2, 3]))
    page, token = store.list(10, "")
    assert [r["id"] for r in page] == [1, 2, 3]
    assert token == "", "nothing left, but next_page_token wasn't empty"

    store = CursorStore(_records([1, 2, 3, 4, 5]))
    page, token = store.list(2, "")
    assert [r["id"] for r in page] == [1, 2]
    assert token != "", "more records remain, but next_page_token was empty"


def test_tampered_token_errors():
    store = CursorStore(_records(range(1, 10)))
    _, token = store.list(3, "")
    with pytest.raises(ValueError):
        store.list(3, _tamper(token))


def test_page_size_clamps():
    store = CursorStore(_records(range(1, 61)))  # 60 records

    page, _ = store.list(0, "")
    assert len(page) == CursorStore.DEFAULT_PAGE_SIZE, "page_size <= 0 didn't fall back to the default"

    page, _ = store.list(-5, "")
    assert len(page) == CursorStore.DEFAULT_PAGE_SIZE, "negative page_size didn't fall back to the default"

    page, _ = store.list(10_000, "")
    assert len(page) == CursorStore.MAX_PAGE_SIZE, "oversized page_size wasn't clamped to the max"


def test_param_change_invalidates_token():
    store = CursorStore(_records(range(1, 30)))
    _, token = store.list(5, "")
    with pytest.raises(ValueError):
        store.list(7, token)


def test_insert_mid_walk_never_duplicates_or_skips():
    seed_ids = [i * 10 for i in range(1, 21)]  # 10, 20, ..., 200
    store = CursorStore(_records(seed_ids))

    page1, token = store.list(5, "")
    assert [r["id"] for r in page1] == seed_ids[:5]

    # Lands before everything returned so far -- the classic case where
    # offset pagination re-shows an already-seen record.
    store.insert({"id": 5, "name": "new-before-cursor"})
    # Lands past the end of the original seed -- must show up later
    # without disturbing anything already walked.
    store.insert({"id": 999, "name": "new-after-cursor"})

    seen = [r["id"] for r in page1]
    while True:
        page, token = store.list(5, token)
        seen.extend(r["id"] for r in page)
        if not token:
            break

    original_seen = [i for i in seen if i in seed_ids]
    assert sorted(original_seen) == seed_ids, "an original record was skipped or duplicated"
    assert len(original_seen) == len(set(original_seen)), "an original record was duplicated"
