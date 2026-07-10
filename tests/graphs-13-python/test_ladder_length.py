from solution import ladder_length


def test_classic():
    word_list = ["hot", "dot", "dog", "lot", "log", "cog"]
    assert ladder_length("hit", "cog", word_list) == 5


def test_end_word_not_in_list():
    word_list = ["hot", "dot", "dog", "lot", "log"]
    assert ladder_length("hit", "cog", word_list) == 0


def test_direct_neighbor():
    assert ladder_length("hit", "hot", ["hot"]) == 2
