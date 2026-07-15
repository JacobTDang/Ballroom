from solution import generate_parenthesis


def test_generate_parenthesis_case_1():
    assert sorted(generate_parenthesis(1)) == sorted(["()"])


def test_generate_parenthesis_case_2():
    assert sorted(generate_parenthesis(2)) == sorted(["(())", "()()"])


def test_generate_parenthesis_case_3():
    assert sorted(generate_parenthesis(3)) == sorted(
        ["((()))", "(()())", "(())()", "()(())", "()()()"]
    )


def test_generate_parenthesis_case_4():
    assert sorted(generate_parenthesis(4)) == sorted(
        [
            "(((())))", "((()()))", "((())())", "((()))()", "(()(()))",
            "(()()())", "(()())()", "(())(())", "(())()()", "()((()))",
            "()(()())", "()(())()", "()()(())", "()()()()",
        ]
    )
