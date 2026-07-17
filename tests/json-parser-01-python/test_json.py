import pytest

from solution import parse


def test_nested_document_exact_structure():
    doc = '{"name": "ada", "age": -3, "tags": ["a", "b"], "meta": {"ok": true, "note": null}, "empty": [], "eo": {}}'
    assert parse(doc) == {
        "name": "ada",
        "age": -3,
        "tags": ["a", "b"],
        "meta": {"ok": True, "note": None},
        "empty": [],
        "eo": {},
    }


def test_escapes():
    assert parse('"say \\"hi\\" and \\\\"') == 'say "hi" and \\'


def test_whitespace_everywhere():
    assert parse('  { "a" :  [ 1 , 2 ]  }  ') == {"a": [1, 2]}


@pytest.mark.parametrize("doc,pos", [
    ('{"a" 1}', "5"),
    ('"unterminated', "0"),
    ("tru", "0"),
    ('{"a": 1} extra', "9"),
    ('"bad \\x escape"', "5"),
])
def test_errors_name_positions(doc, pos):
    with pytest.raises(ValueError, match=pos):
        parse(doc)
