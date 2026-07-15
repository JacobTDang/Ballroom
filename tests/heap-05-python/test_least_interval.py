from solution import least_interval


def test_least_interval_case_1():
    assert least_interval(list("AAABBB"), 2) == 8


def test_least_interval_case_2():
    assert least_interval(list("AAABBB"), 0) == 6


def test_least_interval_case_3():
    assert least_interval(list("AAAAAABCDEFG"), 2) == 16


def test_least_interval_case_4():
    assert least_interval(list("A"), 5) == 1


def test_least_interval_case_5():
    assert least_interval(list("AAAB"), 3) == 9


def test_least_interval_case_6():
    assert least_interval(list("AB"), 2) == 2
