def align(t: int, k: int) -> int:
    """Rounds t down to the start of its k-wide bucket."""
    return (t // k) * k
