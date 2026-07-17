from solution import TTLCache


def test_basic_put_get():
    c = TTLCache(2, 1000)
    c.put_at("a", 1, 0)
    assert c.get_at("a", 10) == 1
    assert c.get_at("missing", 10) is None


def test_lru_eviction_respects_recency():
    c = TTLCache(2, 100000)
    c.put_at("a", 1, 0)
    c.put_at("b", 2, 1)
    c.get_at("a", 2)
    c.put_at("c", 3, 3)
    assert c.get_at("b", 4) is None, "b survived despite being LRU"
    assert c.get_at("a", 4) == 1, "a evicted despite being recently used"
    assert c.get_at("c", 4) == 3


def test_ttl_expiry_from_write_time():
    c = TTLCache(2, 100)
    c.put_at("a", 1, 0)
    assert c.get_at("a", 99) == 1
    assert c.get_at("a", 100) is None, "alive at exactly ttl"


def test_get_refreshes_recency_not_expiry():
    c = TTLCache(2, 100)
    c.put_at("a", 1, 0)
    assert c.get_at("a", 99) == 1
    assert c.get_at("a", 100) is None, "get must not extend the TTL"


def test_expired_entries_do_not_occupy_capacity():
    c = TTLCache(2, 100)
    c.put_at("a", 1, 0)
    c.put_at("b", 2, 0)
    c.put_at("x", 10, 200)
    c.put_at("y", 20, 201)
    assert c.get_at("x", 202) == 10, "a corpse was counted against capacity"
    assert c.get_at("y", 202) == 20, "a corpse was counted against capacity"


def test_rewrite_resets_value_and_expiry():
    c = TTLCache(2, 100)
    c.put_at("a", 1, 0)
    c.put_at("a", 2, 50)
    assert c.get_at("a", 149) == 2
    assert c.get_at("a", 150) is None
