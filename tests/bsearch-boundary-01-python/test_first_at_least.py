from solution import first_at_least


def test_first_at_least_middle():
    assert first_at_least([1, 3, 5, 7, 9], 6) == 3


def test_first_at_least_first_element():
    assert first_at_least([1, 3, 5, 7, 9], 1) == 0


def test_first_at_least_smaller_than_all():
    assert first_at_least([1, 3, 5, 7, 9], 0) == 0


def test_first_at_least_exact_last_element():
    assert first_at_least([1, 3, 5, 7, 9], 9) == 4


def test_first_at_least_past_the_end():
    assert first_at_least([1, 3, 5, 7, 9], 10) == 5


def test_first_at_least_single_element_past_the_end():
    assert first_at_least([5], 10) == 1


def test_first_at_least_duplicates():
    assert first_at_least([2, 2, 2, 2], 2) == 0
