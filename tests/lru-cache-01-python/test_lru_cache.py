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
