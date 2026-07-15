from solution import top_k_frequent


def test_top_k_frequent_case_1():
    assert sorted(top_k_frequent([1, 1, 1, 2, 2, 3], 2)) == [1, 2]


def test_top_k_frequent_case_2():
    assert sorted(top_k_frequent([1], 1)) == [1]


def test_top_k_frequent_case_3():
    assert sorted(top_k_frequent([4, 1, -1, 2, -1, 2, 3], 2)) == [-1, 2]


def test_top_k_frequent_case_4():
    assert sorted(top_k_frequent([5, 5, 5, 5, 3, 3, 1], 1)) == [5]


def test_top_k_frequent_case_5():
    assert sorted(top_k_frequent([1, 2, 3], 3)) == [1, 2, 3]


def test_top_k_frequent_case_6():
    assert sorted(top_k_frequent([1, 1, 1, 1, 2, 2, 2, 3, 3, 4], 2)) == [1, 2]


def test_top_k_frequent_case_7():
    assert sorted(top_k_frequent([-5, -5, -3, -3, -3, -1], 1)) == [-3]


def test_top_k_frequent_case_8():
    assert sorted(top_k_frequent([7, 7, 7], 1)) == [7]


def test_top_k_frequent_case_9():
    assert sorted(top_k_frequent([-10000, -10000, 10000], 1)) == [-10000]
