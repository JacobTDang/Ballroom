def mask(accounts: list[str]) -> list[str]:
    """Return the accounts with all but the last 4 characters starred out."""
    return [a if len(a) < 5 else "*" * (len(a) - 4) + a[-4:] for a in accounts]
