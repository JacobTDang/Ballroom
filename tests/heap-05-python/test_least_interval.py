from solution import least_interval


def test_least_interval():
    assert least_interval(list("AAABBB"), 2) == 8
    assert least_interval(list("AAABBB"), 0) == 6
    assert least_interval(list("AAAAAABCDEFG"), 2) == 16
    assert least_interval(list("A"), 5) == 1
    assert least_interval(list("AAAB"), 3) == 9
    assert least_interval(list("AB"), 2) == 2
