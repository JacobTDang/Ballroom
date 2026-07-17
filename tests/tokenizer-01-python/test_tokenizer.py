import pytest

from solution import tokenize


def test_dense_expression_exact_tokens():
    assert tokenize("3+4.5*x") == [
        ("number", "3", 0),
        ("op", "+", 1),
        ("number", "4.5", 2),
        ("op", "*", 5),
        ("ident", "x", 6),
    ]


def test_parens_and_identifiers():
    assert tokenize("price * (1 + tax_rate2)") == [
        ("ident", "price", 0),
        ("op", "*", 6),
        ("lparen", "(", 8),
        ("number", "1", 9),
        ("op", "+", 11),
        ("ident", "tax_rate2", 13),
        ("rparen", ")", 22),
    ]


def test_second_decimal_point_is_an_error():
    with pytest.raises(ValueError, match="3"):
        tokenize("12..3")


def test_unknown_character_error_names_position():
    with pytest.raises(ValueError, match="2"):
        tokenize("a @ b")


def test_empty_input():
    assert tokenize("") == []
