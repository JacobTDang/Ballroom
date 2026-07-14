from solution import generate_parenthesis


def test_generate_parenthesis():
    assert sorted(generate_parenthesis(1)) == sorted(["()"])
    assert sorted(generate_parenthesis(2)) == sorted(["(())", "()()"])
    assert sorted(generate_parenthesis(3)) == sorted(
        ["((()))", "(()())", "(())()", "()(())", "()()()"]
    )
    assert sorted(generate_parenthesis(4)) == sorted(
        [
            "(((())))", "((()()))", "((())())", "((()))()", "(()(()))",
            "(()()())", "(()())()", "(())(())", "(())()()", "()((()))",
            "()(()())", "()(())()", "()()(())", "()()()()",
        ]
    )
