from solution import find_cheapest_price


def test_one_stop():
    flights = [[0, 1, 100], [1, 2, 100], [2, 0, 100], [1, 3, 600], [2, 3, 200]]
    assert find_cheapest_price(4, flights, 0, 3, 1) == 700


def test_cheaper_via_stop():
    flights = [[0, 1, 100], [1, 2, 100], [0, 2, 500]]
    assert find_cheapest_price(3, flights, 0, 2, 1) == 200


def test_no_stops_allowed():
    flights = [[0, 1, 100], [1, 2, 100], [0, 2, 500]]
    assert find_cheapest_price(3, flights, 0, 2, 0) == 500
