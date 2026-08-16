from solution import next_higher


def test_example_rising_then_drop():
    assert next_higher([30, 40, 50, 20]) == [1, 1, 0, 0]


def test_example_classic_sequence():
    spend = [73, 74, 75, 71, 69, 72, 76, 73]
    assert next_higher(spend) == [1, 1, 4, 2, 1, 1, 0, 0]


def test_example_strictly_decreasing():
    assert next_higher([50, 40, 30]) == [0, 0, 0]


def test_empty():
    assert next_higher([]) == []


def test_single_day():
    assert next_higher([7]) == [0]


def test_all_equal_never_strictly_higher():
    assert next_higher([5, 5, 5]) == [0, 0, 0]


def test_valley():
    assert next_higher([10, 5, 20]) == [2, 1, 0]


def test_plateau_then_rise():
    assert next_higher([4, 4, 5]) == [2, 1, 0]
