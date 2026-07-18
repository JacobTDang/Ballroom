# Accumulates formatted report lines. It's meant to be scratch space
# for one call to generate_report, not storage that outlives it.
_report_lines: list[str] = []


def generate_report(items: list[str]) -> list[str]:
    """Format items (low-stock item names) into report lines, one per
    item. Currently a call's report can contain lines left over from
    an earlier call -- find and fix the bug."""
    for item in items:
        _report_lines.append(f"LOW STOCK: {item}")
    return _report_lines
