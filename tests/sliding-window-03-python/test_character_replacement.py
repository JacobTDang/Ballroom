from solution import character_replacement


def test_character_replacement_case_1():
    assert character_replacement("ABAB", 2) == 4


def test_character_replacement_case_2():
    assert character_replacement("AABABBA", 1) == 4


def test_character_replacement_case_3():
    assert character_replacement("ABCDE", 1) == 2


def test_character_replacement_case_4():
    assert character_replacement("AAAA", 0) == 4


def test_character_replacement_case_5():
    assert character_replacement("A", 0) == 1


def test_character_replacement_case_6():
    assert character_replacement("ABBB", 2) == 4


def test_character_replacement_case_7():
    assert character_replacement("", 2) == 0


def test_character_replacement_case_8():
    assert character_replacement("AAAA", 4) == 4


def test_character_replacement_case_9():
    assert character_replacement("ABABABAB", 3) == 7
