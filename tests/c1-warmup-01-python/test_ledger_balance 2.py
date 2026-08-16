from solution import final_balance


def test_example_mixed():
    assert final_balance(["D 100", "W 30", "W 90"]) == (70, 1)


def test_declined_only():
    assert final_balance(["W 10"]) == (0, 1)


def test_empty():
    assert final_balance([]) == (0, 0)


def test_exact_zero_allowed():
    assert final_balance(["D 50", "W 50"]) == (0, 0)


def test_declines_do_not_change_balance():
    assert final_balance(["D 10", "W 20", "W 20", "D 5"]) == (15, 2)


def test_many_deposits():
    assert final_balance(["D 1"] * 1000) == (1000, 0)
