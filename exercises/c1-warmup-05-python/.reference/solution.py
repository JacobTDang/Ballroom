def is_alternating(digits: str) -> bool:
    """Return True if digit parity strictly alternates along the string."""
    return all((int(a) + int(b)) % 2 == 1 for a, b in zip(digits, digits[1:]))
