from solution import top_k


def test_example_counts_beat_order():
    events = ["acct9", "acct1", "acct9", "acct4", "acct1", "acct9"]
    assert top_k(events, 2) == ["acct9", "acct1"]


def test_example_all_tied_lexicographic():
    assert top_k(["b", "a", "c"], 2) == ["a", "b"]


def test_example_single_event():
    assert top_k(["x"], 1) == ["x"]


def test_k_equals_distinct_count():
    assert top_k(["a", "b", "b"], 2) == ["b", "a"]


def test_tie_broken_lexicographically_within_counts():
    assert top_k(["z", "z", "m", "m", "a"], 2) == ["m", "z"]


def test_empty_events_k_zero():
    assert top_k([], 0) == []


def test_one_account_repeated():
    assert top_k(["a", "a", "a"], 1) == ["a"]


def test_mixed_counts_and_ties():
    events = ["a"] * 3 + ["b"] * 3 + ["c"] * 2 + ["d"]
    assert top_k(events, 3) == ["a", "b", "c"]
