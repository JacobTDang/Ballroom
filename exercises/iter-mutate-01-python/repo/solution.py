def remove_value(values: list[int], target: int) -> list[int]:
    """Remove every occurrence of target from values, in place, and
    return it. Currently leaves some matches behind — find and fix the
    bug."""
    i = 0
    while i < len(values):
        if values[i] == target:
            values.pop(i)
        i += 1
    return values
