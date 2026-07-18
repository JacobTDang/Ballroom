def build_log(chunks: list[str]) -> str:
    """Join chunks (newest-first) into a single oldest-first log, no
    separator."""
    return "".join(reversed(chunks))
