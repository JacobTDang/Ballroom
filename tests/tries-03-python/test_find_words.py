from solution import find_words


def check(board, words, want):
    got = find_words([row[:] for row in board], words)
    assert sorted(got) == sorted(want)


def test_find_words_case_1():
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


def test_find_words_case_2():
    check([list("ab"), list("cd")], ["abcb"], [])


def test_find_words_case_3():
    check([list("a")], ["a"], ["a"])


def test_find_words_case_4():
    check([list("aa")], ["aaa"], [])


def test_find_words_case_5():
    check([list("ab"), list("cd")], ["abdc"], ["abdc"])
