from solution import longest_consecutive


def test_longest_consecutive():
    assert longest_consecutive([100, 4, 200, 1, 3, 2]) == 4
    assert longest_consecutive([0, 3, 7, 2, 5, 8, 4, 6, 0, 1]) == 9
    assert longest_consecutive([]) == 0
    assert longest_consecutive([1, 2, 0, 1]) == 3
    assert longest_consecutive([9, 1, 4, 7, 3, -1, 0, 5, 8, -1, 6]) == 7
    assert longest_consecutive([5]) == 1
    assert longest_consecutive([7, 7, 7, 7]) == 1
    assert longest_consecutive([1, 2, 3, 10, 11]) == 3
    assert longest_consecutive([-3, -2, -1, 0, 1]) == 5
    assert longest_consecutive([50, 3, 51, 2, 52, 1, 4, 49, 48, 47]) == 6
    assert longest_consecutive([-1000000000, -999999999, -999999998]) == 3
