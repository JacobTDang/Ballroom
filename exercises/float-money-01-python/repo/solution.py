def settles_bill(amounts: list[float], bill: float) -> bool:
    """Returns whether amounts sums to bill. Currently unreliable --
    find and fix the bug."""
    total = 0.0
    for amount in amounts:
        total += amount
    return total == bill
