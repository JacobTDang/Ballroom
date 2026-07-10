from solution import car_fleet


def test_car_fleet():
    assert car_fleet(12, [10, 8, 0, 5, 3], [2, 4, 1, 1, 3]) == 3
    assert car_fleet(10, [3], [3]) == 1
    assert car_fleet(100, [0, 2, 4], [4, 2, 1]) == 1
    assert car_fleet(10, [0, 4, 8], [1, 1, 1]) == 3
