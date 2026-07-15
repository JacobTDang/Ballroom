from solution import ParkingLot


def test_initial_availability_matches_construction():
    lot = ParkingLot(2, 3, 4)
    assert lot.available("motorcycle") == 2
    assert lot.available("compact") == 3
    assert lot.available("large") == 4


def test_motorcycle_takes_motorcycle_spot_first():
    lot = ParkingLot(1, 1, 1)
    assert lot.park("motorcycle") > 0
    assert lot.available("motorcycle") == 0
    assert lot.available("compact") == 1


def test_motorcycle_overflows_to_compact_then_large():
    lot = ParkingLot(1, 1, 1)
    lot.park("motorcycle")
    lot.park("motorcycle")
    assert lot.available("compact") == 0
    lot.park("motorcycle")
    assert lot.available("large") == 0


def test_car_takes_compact_before_large():
    lot = ParkingLot(1, 1, 1)
    assert lot.park("car") > 0
    assert lot.available("compact") == 0
    assert lot.available("large") == 1


def test_car_never_takes_motorcycle_spot():
    lot = ParkingLot(5, 0, 0)
    assert lot.park("car") == -1
    assert lot.available("motorcycle") == 5


def test_bus_needs_a_large_spot():
    lot = ParkingLot(5, 5, 0)
    assert lot.park("bus") == -1
    lot2 = ParkingLot(0, 0, 1)
    assert lot2.park("bus") > 0


def test_park_returns_minus_one_when_full():
    lot = ParkingLot(0, 1, 0)
    assert lot.park("car") > 0
    assert lot.park("car") == -1


def test_tickets_are_unique_and_increasing():
    lot = ParkingLot(0, 3, 0)
    t1 = lot.park("car")
    t2 = lot.park("car")
    t3 = lot.park("car")
    assert t1 == 1 and t2 == 2 and t3 == 3


def test_leave_frees_the_right_spot_type():
    lot = ParkingLot(0, 0, 1)
    t = lot.park("bus")
    assert lot.available("large") == 0
    assert lot.leave(t) is True
    assert lot.available("large") == 1
    assert lot.park("bus") > 0


def test_leave_unknown_ticket_returns_false():
    lot = ParkingLot(1, 1, 1)
    assert lot.leave(99) is False


def test_leave_twice_returns_false():
    lot = ParkingLot(0, 1, 0)
    t = lot.park("car")
    assert lot.leave(t) is True
    assert lot.leave(t) is False
    assert lot.available("compact") == 1


def test_tickets_are_not_reused_after_leave():
    lot = ParkingLot(0, 1, 0)
    t1 = lot.park("car")
    lot.leave(t1)
    t2 = lot.park("car")
    assert t2 != t1
