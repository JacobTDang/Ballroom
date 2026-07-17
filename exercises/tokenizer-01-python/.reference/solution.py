def tokenize(input):
    """One pass, one branch per character class -- each consumes a
    maximal token and records where it started. Errors name the exact
    position, because "invalid input" without one is a debugging
    session."""
    tokens = []
    i = 0
    n = len(input)
    while i < n:
        c = input[i]
        if c in " \t\n":
            i += 1
        elif c.isdigit():
            start = i
            saw_dot = False
            while i < n and (input[i].isdigit() or input[i] == "."):
                if input[i] == ".":
                    if saw_dot:
                        raise ValueError(f"second decimal point at position {i}")
                    saw_dot = True
                i += 1
            tokens.append(("number", input[start:i], start))
        elif c.isalpha() or c == "_":
            start = i
            while i < n and (input[i].isalnum() or input[i] == "_"):
                i += 1
            tokens.append(("ident", input[start:i], start))
        elif c in "+-*/":
            tokens.append(("op", c, i))
            i += 1
        elif c == "(":
            tokens.append(("lparen", c, i))
            i += 1
        elif c == ")":
            tokens.append(("rparen", c, i))
            i += 1
        else:
            raise ValueError(f"unexpected character {c!r} at position {i}")
    return tokens
