import pytest

from solution import decode, encode

CASES = [
    ["neet", "code", "love", "you"],
    [],
    [""],
    ["", "", ""],
    ["a#b", "c##d", "5#hello"],
    ["hello world", "foo,bar", "123"],
    ["4#abcd", "hello"],
    ["#####"],
    ["xyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxy"],
    ["123", "456", "0"],
    ["", "a", "", "b"],
]


@pytest.mark.parametrize("strs", CASES)
def test_encode_decode_round_trip(strs):
    encoded = encode(strs)
    assert decode(encoded) == strs
