# First Repeated Merchant

A card's transaction feed lists merchant names in the order the charges
arrived. Duplicate-charge detection cares about the first merchant that
shows up a second time.

Write `first_repeat(merchants: list[str]) -> str | None` returning the
merchant whose second occurrence comes earliest in the list, or None if
no merchant ever repeats.

## Examples

    first_repeat(["uber", "cafe", "uber", "cafe"])  -> "uber"
    first_repeat(["a", "b", "c"])                   -> None
    first_repeat([])                                -> None

## Constraints

- 0 <= len(merchants) <= 10_000
- 1 <= len(merchant) <= 30
