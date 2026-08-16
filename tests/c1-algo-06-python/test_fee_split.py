from solution import min_transfers_two


def test_example_one_creditor_two_debtors():
    assert min_transfers_two([4, -2, -2]) == 2


def test_example_three_debtors():
    assert min_transfers_two([3, -1, -1, -1]) == 3


def test_example_all_zero():
    assert min_transfers_two([0, 0]) == 0


def test_empty():
    assert min_transfers_two([]) == 0


def test_single_pair():
    assert min_transfers_two([5, -5]) == 1


def test_tied_picks_use_lower_index():
    assert min_transfers_two([2, -1, -1, 1, -1]) == 3


def test_two_matched_pairs():
    assert min_transfers_two([1, -1, 1, -1]) == 2


def test_zeros_are_ignored():
    assert min_transfers_two([0, 3, 0, -3, 0]) == 1


def test_creditor_drained_across_transfers():
    assert min_transfers_two([7, -4, -3]) == 2
