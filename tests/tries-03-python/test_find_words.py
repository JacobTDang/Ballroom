from solution import find_words


def check(board, words, want):
    got = find_words([row[:] for row in board], words)
    assert sorted(got) == sorted(want)


def test_find_words():
    check(
        [
            list("oaan"),
            list("etae"),
            list("ihkr"),
            list("iflv"),
        ],
        ["oath", "pea", "eat", "rain"],
        ["eat", "oath"],
    )
    check([list("ab"), list("cd")], ["abcb"], [])
    check([list("a")], ["a"], ["a"])
