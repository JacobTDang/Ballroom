from solution import KthLargest


def test_kth_largest():
    kl = KthLargest(3, [4, 5, 8, 2])
    assert kl.add(3) == 4
    assert kl.add(5) == 5
    assert kl.add(10) == 5
    assert kl.add(9) == 8
    assert kl.add(4) == 8


def test_kth_largest_empty_initial_stream():
    kl = KthLargest(1, [])
    assert kl.add(-3) == -3
    assert kl.add(-2) == -2


def test_kth_largest_k_equals_two():
    kl = KthLargest(2, [0])
    assert kl.add(-1) == -1
    assert kl.add(1) == 0
