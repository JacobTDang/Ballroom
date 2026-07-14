from solution import MedianFinder


def test_median_finder():
    mf = MedianFinder()
    mf.add_num(1)
    mf.add_num(2)
    assert mf.find_median() == 1.5
    mf.add_num(3)
    assert mf.find_median() == 2.0


def test_median_finder_single_element():
    mf = MedianFinder()
    mf.add_num(42)
    assert mf.find_median() == 42.0


def test_median_finder_out_of_order_inserts():
    mf = MedianFinder()
    for n in [5, 1, 9, 3, 7]:
        mf.add_num(n)
    assert mf.find_median() == 5.0
    mf.add_num(10)
    assert mf.find_median() == 6.0


def test_median_finder_negative_values():
    mf = MedianFinder()
    for n in [-5, -1, -3]:
        mf.add_num(n)
    assert mf.find_median() == -3.0
    mf.add_num(-2)
    assert mf.find_median() == -2.5
