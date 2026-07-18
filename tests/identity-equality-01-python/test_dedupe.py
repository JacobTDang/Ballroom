from solution import Record, dedupe


def _kv(records):
    return [(r.key, r.value) for r in records]


def test_value_equal_but_distinct_objects_collapse():
    r1 = Record("a", 1)
    r2 = Record("a", 1)  # distinct object, same value
    r3 = Record("b", 2)
    assert _kv(dedupe([r1, r2, r3])) == [("a", 1), ("b", 2)]


def test_triple_duplicate_collapses_to_one():
    r1 = Record("x", 5)
    r2 = Record("x", 5)
    r3 = Record("x", 5)
    assert _kv(dedupe([r1, r2, r3])) == [("x", 5)]


def test_all_distinct_records_survive():
    r1 = Record("a", 1)
    r2 = Record("b", 2)
    r3 = Record("c", 3)
    assert _kv(dedupe([r1, r2, r3])) == [("a", 1), ("b", 2), ("c", 3)]


def test_interleaved_duplicates_preserve_first_occurrence_order():
    r1 = Record("a", 1)
    r2 = Record("b", 2)
    r3 = Record("a", 1)  # dup of r1
    r4 = Record("c", 3)
    r5 = Record("b", 2)  # dup of r2
    assert _kv(dedupe([r1, r2, r3, r4, r5])) == [("a", 1), ("b", 2), ("c", 3)]


def test_same_key_different_value_is_not_a_duplicate():
    r1 = Record("a", 1)
    r2 = Record("a", 2)
    assert _kv(dedupe([r1, r2])) == [("a", 1), ("a", 2)]
