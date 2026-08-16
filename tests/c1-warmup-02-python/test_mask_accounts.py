from solution import mask


def test_example_long_number():
    assert mask(["123456789"]) == ["*****6789"]


def test_example_mixed_lengths():
    assert mask(["9604", "10011001"]) == ["9604", "****1001"]


def test_empty_list():
    assert mask([]) == []


def test_exactly_four_unchanged():
    assert mask(["1234"]) == ["1234"]


def test_exactly_five_masks_one():
    assert mask(["12345"]) == ["*2345"]


def test_single_char_unchanged():
    assert mask(["7"]) == ["7"]


def test_order_preserved():
    assert mask(["11111", "22", "3333333"]) == ["*1111", "22", "***3333"]


def test_long_account():
    assert mask(["1" * 16]) == ["*" * 12 + "1111"]
