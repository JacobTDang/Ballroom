from solution import eval_rpn


def test_eval_rpn_case_1():
    assert eval_rpn(["2", "1", "+", "3", "*"]) == 9


def test_eval_rpn_case_2():
    assert eval_rpn(["4", "13", "5", "/", "+"]) == 6


def test_eval_rpn_case_3():
    assert (
        eval_rpn(["10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"])
        == 22
    )


def test_eval_rpn_case_4():
    assert eval_rpn(["18"]) == 18


def test_eval_rpn_case_5():
    assert eval_rpn(["4", "3", "-"]) == 1


def test_eval_rpn_case_6():
    assert eval_rpn(["-3", "4", "+"]) == 1


def test_eval_rpn_case_7():
    assert eval_rpn(["7", "-2", "/"]) == -3


def test_eval_rpn_case_8():
    assert eval_rpn(["5", "5", "*", "5", "*"]) == 125
