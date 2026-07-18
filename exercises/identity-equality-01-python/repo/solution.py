class Record:
    def __init__(self, key: str, value: int):
        self.key = key
        self.value = value

    def __repr__(self):
        return f"Record({self.key!r}, {self.value})"


def dedupe(records: list[Record]) -> list[Record]:
    """Removes duplicate records -- two records with the same key and
    value are duplicates. Currently keeps both -- find and fix the
    bug."""
    seen = set()
    result = []
    for r in records:
        if r not in seen:
            seen.add(r)
            result.append(r)
    return result
