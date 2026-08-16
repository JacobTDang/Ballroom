from solution import reverse_words


def test_example_sentence():
    assert reverse_words("pay the balance now") == "now balance the pay"


def test_example_single_word():
    assert reverse_words("done!") == "done!"


def test_example_two_words():
    assert reverse_words("a b") == "b a"


def test_punctuation_stays_with_word():
    assert reverse_words("hello, world!") == "world! hello,"


def test_palindrome_order():
    assert reverse_words("x y x") == "x y x"


def test_single_char_word():
    assert reverse_words("q") == "q"


def test_repeated_words():
    assert reverse_words("go go stop") == "stop go go"


def test_longer_sentence():
    assert reverse_words("one two three four five") == "five four three two one"
