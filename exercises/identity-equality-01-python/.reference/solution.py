class Record:
    def __init__(self, key: str, value: int):
        self.key = key
        self.value = value

    def __repr__(self):
        return f"Record({self.key!r}, {self.value})"

    def __eq__(self, other):
        return isinstance(other, Record) and self.key == other.key and self.value == other.value

    def __hash__(self):
        return hash((self.key, self.value))


def dedupe(records: list[Record]) -> list[Record]:
    """Removes duplicate records -- two records with the same key and
    value are duplicates."""
    seen = set()
    result = []
    for r in records:
        if r not in seen:
            seen.add(r)
            result.append(r)
    return result
