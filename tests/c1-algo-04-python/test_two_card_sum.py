from solution import pair_within_budget


def test_example_exact_budget_pair():
    assert pair_within_budget([40, 60, 25, 30], 90) == (1, 3)


def test_example_tie_smallest_i_then_j():
    assert pair_within_budget([50, 40, 40, 50], 90) == (0, 1)


def test_example_no_pair_fits():
    assert pair_within_budget([10, 20, 15], 12) is None


def test_all_equal_prices():
    assert pair_within_budget([5, 5, 5], 10) == (0, 1)


def test_two_cards_exact_boundary():
    assert pair_within_budget([3, 7], 10) == (0, 1)


def test_single_price_no_pair():
    assert pair_within_budget([8], 100) is None


def test_empty_prices():
    assert pair_within_budget([], 5) is None


def test_skips_large_card_for_fitting_pair():
    assert pair_within_budget([9, 1, 2], 4) == (1, 2)


def test_duplicate_values_complement_equals_price():
    assert pair_within_budget([4, 4, 4, 2], 8) == (0, 1)


def test_best_total_below_budget():
    assert pair_within_budget([1, 8, 3], 9) == (0, 1)
