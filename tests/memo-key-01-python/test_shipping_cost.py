from solution import shipping_cost


def test_same_weight_different_zones_in_order():
    # Weight 5kg, zones queried A -> B -> C in order: each must be
    # correct even though the previous call just cached weight 5 for a
    # different zone.
    assert shipping_cost(5, "A") == 1500
    assert shipping_cost(5, "B") == 3000
    assert shipping_cost(5, "C") == 5000


def test_same_weight_different_zones_reversed_order():
    # Weight 8kg, zones queried in a different order (C -> A -> B) --
    # the bug must not survive a different call sequence either.
    assert shipping_cost(8, "C") == 7700
    assert shipping_cost(8, "A") == 2100
    assert shipping_cost(8, "B") == 4500
