from solution import daily_temperatures


def test_daily_temperatures():
    assert daily_temperatures([73, 74, 75, 71, 69, 72, 76, 73]) == [1, 1, 4, 2, 1, 1, 0, 0]
    assert daily_temperatures([30, 40, 50, 60]) == [1, 1, 1, 0]
    assert daily_temperatures([30, 60, 90]) == [1, 1, 0]
    assert daily_temperatures([80, 79, 78]) == [0, 0, 0]
    assert daily_temperatures([75]) == [0]
    assert daily_temperatures([55, 55, 55, 60]) == [3, 2, 1, 0]
