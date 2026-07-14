from solution import product_except_self


def test_product_except_self():
    assert product_except_self([1, 2, 3, 4]) == [24, 12, 8, 6]
    assert product_except_self([-1, 1, 0, -3, 3]) == [0, 0, 9, 0, 0]
    assert product_except_self([2, 3]) == [3, 2]
    assert product_except_self([5, 0, 0, 4]) == [0, 0, 0, 0]
    assert product_except_self([1, 1, 1, 1]) == [1, 1, 1, 1]
    assert product_except_self([-1, -2, -3, -4]) == [-24, -12, -8, -6]
    assert product_except_self([1, 2, 3, 4, 5, 6, 7, 8]) == [40320, 20160, 13440, 10080, 8064, 6720, 5760, 5040]
    assert product_except_self([-1, 2, -3, 4]) == [-24, 12, -8, 6]
    assert product_except_self([30, -30, 1]) == [-30, 30, -900]
    assert product_except_self([1, 0, 3, 4]) == [0, 12, 0, 0]
