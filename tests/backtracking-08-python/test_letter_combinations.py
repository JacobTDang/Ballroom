from solution import letter_combinations


def test_letter_combinations():
    assert sorted(letter_combinations("23")) == sorted(
        ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"]
    )
    assert letter_combinations("") == []
    assert sorted(letter_combinations("2")) == sorted(["a", "b", "c"])
    assert sorted(letter_combinations("9")) == sorted(["w", "x", "y", "z"])
    assert sorted(letter_combinations("79")) == sorted(
        ["pw", "px", "py", "pz", "qw", "qx", "qy", "qz", "rw", "rx", "ry", "rz", "sw", "sx", "sy", "sz"]
    )
