def parse(input):
    """Parse a JSON subset: objects, arrays, strings (\\\" and \\\\
    escapes), integers, true/false/null.

    TODO: this handles only a flat {"key": "value"} object via string
    splitting -- no nesting, no arrays, no numbers, no real errors.
    """
    body = input.strip().strip("{}")
    result = {}
    for pair in body.split(","):
        if ":" in pair:
            k, v = pair.split(":", 1)
            result[k.strip().strip('"')] = v.strip().strip('"')
    return result
