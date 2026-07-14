from solution import WordDictionary


def test_word_dictionary():
    d = WordDictionary()
    d.add_word("bad")
    d.add_word("dad")
    d.add_word("mad")

    assert d.search("pad") is False
    assert d.search("bad") is True
    assert d.search(".ad") is True
    assert d.search("b..") is True
    assert d.search("...") is True
    assert d.search("....") is False
    assert d.search("..d") is True
    assert d.search("dab") is False


def test_empty_dictionary_never_matches():
    d = WordDictionary()
    assert d.search("a") is False
    assert d.search(".") is False


def test_wrong_length_queries():
    d = WordDictionary()
    d.add_word("bad")
    d.add_word("dad")
    d.add_word("mad")

    assert d.search(".") is False
    assert d.search("ba.") is True
