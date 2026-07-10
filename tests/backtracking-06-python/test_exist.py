from solution import exist


def make_board():
    return [
        list("ABCE"),
        list("SFCS"),
        list("ADEE"),
    ]


def test_exist():
    assert exist(make_board(), "ABCCED") is True
    assert exist(make_board(), "SEE") is True
    assert exist(make_board(), "ABCB") is False
    assert exist(make_board(), "ABFSAB") is False


def test_exist_single_cell():
    assert exist([["A"]], "A") is True
    assert exist([["A"]], "AA") is False
