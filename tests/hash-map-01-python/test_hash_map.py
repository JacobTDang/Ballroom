from solution import MyHashMap


def test_put_then_get_roundtrip():
    m = MyHashMap()
    m.put(1, 100)
    assert m.get(1) == 100


def test_get_missing_returns_minus_one():
    m = MyHashMap()
    assert m.get(42) == -1


def test_put_overwrites_existing_key():
    m = MyHashMap()
    m.put(1, 100)
    m.put(1, 200)
    assert m.get(1) == 200


def test_remove_then_get_returns_minus_one():
    m = MyHashMap()
    m.put(7, 70)
    m.remove(7)
    assert m.get(7) == -1


def test_remove_missing_key_is_a_noop():
    m = MyHashMap()
    m.put(1, 10)
    m.remove(99)
    assert m.get(1) == 10


def test_colliding_keys_are_kept_separate():
    m = MyHashMap()
    # These collide under any bucket count that divides 1024/1000-style
    # tables; a broken chain loses one of them.
    for k in (1, 1025, 2049, 1001, 2001):
        m.put(k, k * 3)
    for k in (1, 1025, 2049, 1001, 2001):
        assert m.get(k) == k * 3


def test_removing_one_colliding_key_keeps_the_others():
    m = MyHashMap()
    m.put(1, 11)
    m.put(1025, 22)
    m.put(2049, 33)
    m.remove(1025)
    assert m.get(1) == 11
    assert m.get(1025) == -1
    assert m.get(2049) == 33


def test_zero_key_and_zero_value():
    m = MyHashMap()
    m.put(0, 0)
    assert m.get(0) == 0


def test_large_key_bounds():
    m = MyHashMap()
    m.put(1000000, 999)
    assert m.get(1000000) == 999


def test_many_keys_no_aliasing():
    m = MyHashMap()
    for k in range(500):
        m.put(k, k * 2)
    for k in range(500):
        assert m.get(k) == k * 2


def test_interleaved_put_remove_sequence():
    m = MyHashMap()
    m.put(5, 1)
    m.put(6, 2)
    m.remove(5)
    m.put(6, 3)
    m.put(5, 4)
    assert m.get(5) == 4
    assert m.get(6) == 3
