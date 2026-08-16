import heapq


def min_transfers_two(balances: list[int]) -> int:
    """Return the number of transfers made by the settle-up procedure."""
    debtors = [(b, i) for i, b in enumerate(balances) if b < 0]
    creditors = [(-b, i) for i, b in enumerate(balances) if b > 0]
    heapq.heapify(debtors)
    heapq.heapify(creditors)
    transfers = 0
    while debtors:
        debt, debtor = heapq.heappop(debtors)      # debt < 0, largest owed first
        credit, creditor = heapq.heappop(creditors)  # credit < 0, largest owed-to first
        settled = min(-debt, -credit)
        transfers += 1
        debt += settled
        credit += settled
        if debt < 0:
            heapq.heappush(debtors, (debt, debtor))
        if credit < 0:
            heapq.heappush(creditors, (credit, creditor))
    return transfers
