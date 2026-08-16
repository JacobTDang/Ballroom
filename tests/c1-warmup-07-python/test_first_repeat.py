from solution import first_repeat


def test_example_earliest_second_occurrence():
    assert first_repeat(["uber", "cafe", "uber", "cafe"]) == "uber"


def test_example_no_repeats():
    assert first_repeat(["a", "b", "c"]) is None


def test_example_empty():
    assert first_repeat([]) is None


def test_single_element():
    assert first_repeat(["solo"]) is None


def test_immediate_repeat():
    assert first_repeat(["gas", "gas", "cafe"]) == "gas"


def test_later_first_seen_wins_by_second_occurrence():
    assert first_repeat(["a", "b", "b", "a"]) == "b"


def test_repeat_at_end():
    assert first_repeat(["x", "y", "z", "x"]) == "x"


def test_triple_returns_on_second():
    assert first_repeat(["m", "n", "m", "m"]) == "m"
