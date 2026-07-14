from solution import INF, walls_and_gates


def test_walls_and_gates():
    rooms = [
        [INF, -1, 0, INF],
        [INF, INF, INF, -1],
        [INF, -1, INF, -1],
        [0, -1, INF, INF],
    ]
    want = [
        [3, -1, 0, 1],
        [2, 2, 1, -1],
        [1, -1, 2, -1],
        [0, -1, 3, 4],
    ]
    walls_and_gates(rooms)
    assert rooms == want


def test_walls_and_gates_unreachable_room_stays_inf():
    rooms = [[0, -1, INF]]
    want = [[0, -1, INF]]
    walls_and_gates(rooms)
    assert rooms == want


def test_walls_and_gates_no_gates():
    rooms = [[INF, INF]]
    want = [[INF, INF]]
    walls_and_gates(rooms)
    assert rooms == want


def test_walls_and_gates_nearest_gate_wins():
    rooms = [[0, INF, INF, 0]]
    want = [[0, 1, 1, 0]]
    walls_and_gates(rooms)
    assert rooms == want
