from solution import align


def test_align_positive_below_bucket():
    assert align(7, 4) == 4


def test_align_positive_exact_multiple():
    assert align(8, 4) == 8


def test_align_zero():
    assert align(0, 4) == 0


def test_align_negative_near_zero():
    assert align(-1, 4) == -4


def test_align_negative_pinned_case_1():
    assert align(-7, 4) == -8


def test_align_negative_pinned_case_2_exact_multiple():
    assert align(-8, 4) == -8


def test_align_negative_further_out():
    assert align(-9, 4) == -12


def test_align_negative_other_bucket_width():
    assert align(-10, 3) == -12


def test_align_negative_exact_multiple_other_width():
    assert align(-12, 3) == -12
