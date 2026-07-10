from solution import partition


def normalize(lists):
    return sorted(lists)


def test_partition():
    assert normalize(partition("aab")) == normalize([["a", "a", "b"], ["aa", "b"]])
    assert normalize(partition("a")) == [["a"]]
    assert normalize(partition("aba")) == normalize([["a", "b", "a"], ["aba"]])
