from solution import generate_report


def test_calls_do_not_leak_into_each_other():
    # One sequence of three calls exercises all three traps: a lone
    # call passes even with the bug (nothing to leak from yet), a
    # second call with different input is polluted by the first, and
    # a third call repeating the first call's input shows that
    # pollution as an outright duplicate.
    assert generate_report(["apples"]) == ["LOW STOCK: apples"]
    assert generate_report(["bananas"]) == ["LOW STOCK: bananas"]
    assert generate_report(["apples"]) == ["LOW STOCK: apples"]
