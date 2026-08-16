from solution import min_removals


def test_example_one_unclosed_open():
    assert min_removals("(()") == 1


def test_example_close_then_open():
    assert min_removals("())(") == 2


def test_example_already_balanced():
    assert min_removals("()()") == 0


def test_empty():
    assert min_removals("") == 0


def test_all_closes():
    assert min_removals(")))") == 3


def test_all_opens():
    assert min_removals("(((") == 3


def test_reversed_pair():
    assert min_removals(")(") == 2


def test_nested_balanced():
    assert min_removals("(()())") == 0


def test_long_boundary_one_extra_open():
    assert min_removals("(" * 1000 + ")" * 999) == 1
