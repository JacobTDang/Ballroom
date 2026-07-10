def max_profit(prices: list[int]) -> int:
    """Return the maximum profit from buying on one day and selling on
    a later day, or 0 if no profit is possible."""
    if not prices:
        return 0
    min_price = prices[0]
    best = 0
    for p in prices[1:]:
        best = max(best, p - min_price)
        min_price = min(min_price, p)
    return best
