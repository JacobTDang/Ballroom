def settles_bill(amounts: list[float], bill: float) -> bool:
    """Returns whether amounts sums to bill, to the nearest cent."""
    total = 0.0
    for amount in amounts:
        total += amount
    return round(total * 100) == round(bill * 100)
