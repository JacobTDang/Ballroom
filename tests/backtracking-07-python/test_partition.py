from solution import partition


def normalize(lists):
    return sorted(lists)


def test_partition_case_1():
    assert normalize(partition("aab")) == normalize([["a", "a", "b"], ["aa", "b"]])


def test_partition_case_2():
    assert normalize(partition("a")) == [["a"]]


def test_partition_case_3():
    assert normalize(partition("aba")) == normalize([["a", "b", "a"], ["aba"]])


def test_partition_case_4():
    assert normalize(partition("aa")) == normalize([["a", "a"], ["aa"]])


def test_partition_case_5():
    assert normalize(partition("abcba")) == normalize(
        [["a", "b", "c", "b", "a"], ["a", "bcb", "a"], ["abcba"]]
    )
