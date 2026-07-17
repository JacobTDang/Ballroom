def parse(input):
    """Read an INI document into {section: {key: value}}.

    TODO: no sections, no comments, no errors -- every line is split
    on '=' into the "" section, and malformed lines are silently
    skipped.
    """
    result = {"": {}}
    for line in input.split("\n"):
        if "=" in line:
            key, value = line.split("=", 1)
            result[""][key] = value
    return result
