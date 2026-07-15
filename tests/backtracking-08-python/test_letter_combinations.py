from solution import letter_combinations


def test_letter_combinations_case_1():
    assert sorted(letter_combinations("23")) == sorted(
        ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"]
    )


def test_letter_combinations_case_2():
    assert letter_combinations("") == []


def test_letter_combinations_case_3():
    assert sorted(letter_combinations("2")) == sorted(["a", "b", "c"])


def test_letter_combinations_case_4():
    assert sorted(letter_combinations("9")) == sorted(["w", "x", "y", "z"])


def test_letter_combinations_case_5():
    assert sorted(letter_combinations("79")) == sorted(
        ["pw", "px", "py", "pz", "qw", "qx", "qy", "qz", "rw", "rx", "ry", "rz", "sw", "sx", "sy", "sz"]
    )
