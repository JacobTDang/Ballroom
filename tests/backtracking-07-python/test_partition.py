from solution import partition


def normalize(lists):
    return sorted(lists)


def test_partition():
    assert normalize(partition("aab")) == normalize([["a", "a", "b"], ["aa", "b"]])
    assert normalize(partition("a")) == [["a"]]
    assert normalize(partition("aba")) == normalize([["a", "b", "a"], ["aba"]])
    assert normalize(partition("aa")) == normalize([["a", "a"], ["aa"]])
    assert normalize(partition("abcba")) == normalize(
        [["a", "b", "c", "b", "a"], ["a", "bcb", "a"], ["abcba"]]
    )
