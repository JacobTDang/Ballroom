from solution import Grid


def test_snapshot_independent_of_later_mutation():
    g = Grid(3, 2)
    g.set(0, 0, 1)
    g.set(1, 1, 2)
    g.set(2, 0, 3)

    snap = g.snapshot()
    assert snap == [[1, 0], [0, 2], [3, 0]]

    g.set(0, 0, 111)
    g.set(2, 1, 999)

    assert snap == [[1, 0], [0, 2], [3, 0]]
    assert g.get(0, 0) == 111
    assert g.get(2, 1) == 999


def test_multiple_snapshots_are_independent_of_each_other():
    g = Grid(3, 2)
    g.set(0, 0, 1)
    g.set(1, 1, 2)
    g.set(2, 0, 3)

    snap1 = g.snapshot()

    g.set(0, 0, 111)
    g.set(2, 1, 999)

    snap2 = g.snapshot()

    g.set(0, 0, 222)

    assert snap1 == [[1, 0], [0, 2], [3, 0]]
    assert snap2 == [[111, 0], [0, 2], [3, 999]]
    assert g.get(0, 0) == 222


def test_snapshot_covers_every_row_not_just_the_first():
    g = Grid(3, 2)
    for r in range(3):
        for c in range(2):
            g.set(r, c, r * 10 + c)

    snap = g.snapshot()
    g.set(2, 1, -1)  # mutate the LAST row after the snapshot

    assert snap[2][1] == 21  # 2*10+1, unaffected
    assert g.get(2, 1) == -1
