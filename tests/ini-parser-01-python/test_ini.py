import pytest

from solution import parse


def test_full_document():
    doc = """top = level
# a comment
; another comment

[server]
host = localhost
port = 8080
host = example.com

[client]
retries=3
"""
    assert parse(doc) == {
        "": {"top": "level"},
        "server": {"host": "example.com", "port": "8080"},
        "client": {"retries": "3"},
    }


def test_whitespace_trimmed_everywhere():
    assert parse("  spaced key   =   spaced value  ") == {"": {"spaced key": "spaced value"}}


def test_malformed_line_errors_with_line_number():
    with pytest.raises(ValueError, match="2"):
        parse("ok = 1\nnot a valid line\nok2 = 2")


def test_unclosed_section_errors_with_line_number():
    with pytest.raises(ValueError, match="3"):
        parse("[server]\nkey = v\n[broken")
