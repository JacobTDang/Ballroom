def eval_rpn(tokens: list[str]) -> int:
    """Evaluate an arithmetic expression given in Reverse Polish
    Notation and return the result."""
    stack: list[int] = []
    for tok in tokens:
        if tok in ("+", "-", "*", "/"):
            b = stack.pop()
            a = stack.pop()
            if tok == "+":
                res = a + b
            elif tok == "-":
                res = a - b
            elif tok == "*":
                res = a * b
            else:
                res = int(a / b)  # truncate toward zero, not floor
            stack.append(res)
        else:
            stack.append(int(tok))
    return stack[0]
