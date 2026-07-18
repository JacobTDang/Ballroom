def build_log(chunks: list[str]) -> str:
    """Join chunks (newest-first) into a single oldest-first log, no
    separator. Currently far slower than it should be on a large page --
    find and fix the bug."""
    s = ""
    for c in chunks:
        s = c + s
    return s
