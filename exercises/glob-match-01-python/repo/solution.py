def match(pattern, s):
    """Whole-string glob: * (any run), ? (exactly one), [a-c] (one
    from a set/range). An unclosed [ makes the pattern match nothing.

    TODO: this handles only a lone "*" and literal equality.
    """
    if pattern == "*":
        return True
    return pattern == s
