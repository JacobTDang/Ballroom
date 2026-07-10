from solution import decode, encode

CASES = [
    ["neet", "code", "love", "you"],
    [],
    [""],
    ["", "", ""],
    ["a#b", "c##d", "5#hello"],
    ["hello world", "foo,bar", "123"],
]


def test_encode_decode_round_trip():
    for strs in CASES:
        encoded = encode(strs)
        assert decode(encoded) == strs
