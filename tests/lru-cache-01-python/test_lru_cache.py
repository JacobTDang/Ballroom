from solution import LRUCache


def test_lru_cache():
    c = LRUCache(2)
    c.put(1, 100)
    c.put(2, 200)
    assert c.get(1) == 100

    c.put(3, 300)  # evicts 2
    assert c.get(2) == -1
    assert c.get(3) == 300

    c.put(4, 400)  # evicts 1
    assert c.get(1) == -1
    assert c.get(3) == 300
    assert c.get(4) == 400


def test_updating_existing_key_does_not_evict():
    c = LRUCache(2)
    c.put(1, 1)
    c.put(2, 2)
    c.put(1, 10)  # update, not a new insert -- must not evict 2
    assert c.get(2) == 2
    assert c.get(1) == 10


def test_get_refreshes_recency():
    c = LRUCache(2)
    c.put(1, 1)
    c.put(2, 2)
    c.get(1)  # 1 is now most recently used
    c.put(3, 3)  # should evict 2, not 1
    assert c.get(2) == -1
    assert c.get(1) == 1


def test_capacity_one_evicts_immediately():
    c = LRUCache(1)
    c.put(1, 1)
    c.put(2, 2)
    assert c.get(1) == -1
    assert c.get(2) == 2


def test_missing_key_returns_negative_one():
    c = LRUCache(2)
    assert c.get(999) == -1
