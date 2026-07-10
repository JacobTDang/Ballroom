from solution import TimeMap


def test_time_map():
    m = TimeMap()
    m.set("foo", "bar", 1)
    assert m.get("foo", 1) == "bar"
    assert m.get("foo", 3) == "bar"
    m.set("foo", "bar2", 4)
    assert m.get("foo", 4) == "bar2"
    assert m.get("foo", 5) == "bar2"


def test_get_before_any_set_returns_empty():
    m = TimeMap()
    m.set("foo", "bar", 5)
    assert m.get("foo", 1) == ""


def test_unknown_key_returns_empty():
    m = TimeMap()
    assert m.get("missing", 1) == ""
