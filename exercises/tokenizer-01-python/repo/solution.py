def tokenize(input):
    """Split input into (kind, text, pos) tuples.

    TODO: splitting on spaces calls "3+4" one token and loses every
    position -- and nothing is ever an error.
    """
    return [("ident", f, 0) for f in input.split()]
