from solution import top_k_frequent


def test_top_k_frequent():
    assert sorted(top_k_frequent([1, 1, 1, 2, 2, 3], 2)) == [1, 2]
    assert sorted(top_k_frequent([1], 1)) == [1]
    assert sorted(top_k_frequent([4, 1, -1, 2, -1, 2, 3], 2)) == [-1, 2]
    assert sorted(top_k_frequent([5, 5, 5, 5, 3, 3, 1], 1)) == [5]
    assert sorted(top_k_frequent([1, 2, 3], 3)) == [1, 2, 3]
    assert sorted(top_k_frequent([1, 1, 1, 1, 2, 2, 2, 3, 3, 4], 2)) == [1, 2]
    assert sorted(top_k_frequent([-5, -5, -3, -3, -3, -1], 1)) == [-3]
    assert sorted(top_k_frequent([7, 7, 7], 1)) == [7]
    assert sorted(top_k_frequent([-10000, -10000, 10000], 1)) == [-10000]
