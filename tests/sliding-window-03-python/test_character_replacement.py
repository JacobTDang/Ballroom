from solution import character_replacement


def test_character_replacement():
    assert character_replacement("ABAB", 2) == 4
    assert character_replacement("AABABBA", 1) == 4
    assert character_replacement("ABCDE", 1) == 2
    assert character_replacement("AAAA", 0) == 4
    assert character_replacement("A", 0) == 1
    assert character_replacement("ABBB", 2) == 4
    assert character_replacement("", 2) == 0
    assert character_replacement("AAAA", 4) == 4
    assert character_replacement("ABABABAB", 3) == 7
