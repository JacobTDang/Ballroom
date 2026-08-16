def final_balance(transactions: list[str]) -> tuple[int, int]:
    """Return (end balance, declined withdrawal count) for the day."""
    balance = declined = 0
    for t in transactions:
        kind, raw = t.split()
        amount = int(raw)
        if kind == "D":
            balance += amount
        elif balance >= amount:
            balance -= amount
        else:
            declined += 1
    return balance, declined
