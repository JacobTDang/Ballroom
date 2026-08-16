def compress(s: str) -> str:
    """Return the run-length encoding of s, or s if that is not shorter."""
    if not s:
        return s
    parts = []
    run_char, run_len = s[0], 1
    for ch in s[1:]:
        if ch == run_char:
            run_len += 1
        else:
            parts.append(run_char if run_len == 1 else f"{run_char}{run_len}")
            run_char, run_len = ch, 1
    parts.append(run_char if run_len == 1 else f"{run_char}{run_len}")
    encoded = "".join(parts)
    return encoded if len(encoded) < len(s) else s
