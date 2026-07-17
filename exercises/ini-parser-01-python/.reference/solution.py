def parse(input):
    """Line-oriented with a current-section cursor. Every branch
    either consumes the line's whole meaning or raises with its
    1-based number -- a config line that parses as nothing is a typo
    the user deserves to hear about."""
    result = {"": {}}
    section = ""
    for n, raw in enumerate(input.split("\n"), start=1):
        line = raw.strip()
        if not line or line.startswith("#") or line.startswith(";"):
            continue
        if line.startswith("["):
            if not line.endswith("]"):
                raise ValueError(f"line {n}: unclosed section header {line!r}")
            section = line[1:-1].strip()
            result.setdefault(section, {})
        elif "=" in line:
            key, value = line.split("=", 1)
            key = key.strip()
            if not key:
                raise ValueError(f"line {n}: empty key")
            result[section][key] = value.strip()  # later key wins
        else:
            raise ValueError(f"line {n}: not a header, comment, or key=value: {line!r}")
    return result
