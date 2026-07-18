def generate_report(items: list[str]) -> list[str]:
    """Format items (low-stock item names) into report lines, one per
    item."""
    report_lines = []
    for item in items:
        report_lines.append(f"LOW STOCK: {item}")
    return report_lines
