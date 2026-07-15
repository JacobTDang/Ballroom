from solution import min_window


def test_min_window_case_1():
    assert min_window("ADOBECODEBANC", "ABC") == "BANC"


def test_min_window_case_2():
    assert min_window("a", "a") == "a"


def test_min_window_case_3():
    assert min_window("a", "aa") == ""


def test_min_window_case_4():
    assert min_window("ab", "b") == "b"


def test_min_window_case_5():
    assert min_window("bba", "ab") == "ba"


def test_min_window_case_6():
    assert min_window("abc", "abc") == "abc"


def test_min_window_case_7():
    assert min_window("aaflslflsldkalskaaa", "aaa") == "aaa"


def test_min_window_case_8():
    assert min_window("cabwefgewcwaefgcf", "cae") == "cwae"
