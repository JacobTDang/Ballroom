from solution import min_window


def test_min_window():
    assert min_window("ADOBECODEBANC", "ABC") == "BANC"
    assert min_window("a", "a") == "a"
    assert min_window("a", "aa") == ""
    assert min_window("ab", "b") == "b"
    assert min_window("bba", "ab") == "ba"
    assert min_window("abc", "abc") == "abc"
    assert min_window("aaflslflsldkalskaaa", "aaa") == "aaa"
    assert min_window("cabwefgewcwaefgcf", "cae") == "cwae"
