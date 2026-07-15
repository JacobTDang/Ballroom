from solution import car_fleet


def test_car_fleet_case_1():
    assert car_fleet(12, [10, 8, 0, 5, 3], [2, 4, 1, 1, 3]) == 3


def test_car_fleet_case_2():
    assert car_fleet(10, [3], [3]) == 1


def test_car_fleet_case_3():
    assert car_fleet(100, [0, 2, 4], [4, 2, 1]) == 1


def test_car_fleet_case_4():
    assert car_fleet(10, [0, 4, 8], [1, 1, 1]) == 3


def test_car_fleet_case_5():
    assert car_fleet(10, [0, 3, 6], [5, 5, 5]) == 3


def test_car_fleet_case_6():
    assert car_fleet(20, [1, 4], [2, 1]) == 1


def test_car_fleet_case_7():
    assert car_fleet(5, [5], [1]) == 1
