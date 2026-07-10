from solution import letter_combinations


def test_letter_combinations():
    assert sorted(letter_combinations("23")) == sorted(
        ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"]
    )
    assert letter_combinations("") == []
    assert sorted(letter_combinations("2")) == sorted(["a", "b", "c"])
