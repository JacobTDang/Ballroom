from functools import cmp_to_key


class Entry:
    def __init__(self, name: str, score: int):
        self.name = name
        self.score = score

    def __repr__(self):
        return f"Entry({self.name!r}, {self.score})"


def _compare(a: Entry, b: Entry) -> int:
    if a.score != b.score:
        return b.score - a.score  # higher score first
    if a.name < b.name:
        return 1
    if a.name > b.name:
        return -1
    return 0


def sort_leaderboard(entries: list[Entry]) -> list[Entry]:
    """Sorts entries by score descending; ties break by name ascending.
    Currently the tie-break is backwards -- find and fix the bug."""
    return sorted(entries, key=cmp_to_key(_compare))
