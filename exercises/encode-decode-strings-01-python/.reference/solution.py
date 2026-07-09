def encode(strs: list[str]) -> str:
    """Encode a list of strings into a single string that decode can
    reconstruct exactly, including strings that contain any characters
    (digits, delimiters, etc). Each string is prefixed with its length
    and a '#' delimiter, so the delimiter itself can safely appear
    inside a string without ambiguity."""
    return "".join(f"{len(s)}#{s}" for s in strs)


def decode(s: str) -> list[str]:
    """Reverse encode."""
    result = []
    i = 0
    while i < len(s):
        j = i
        while s[j] != "#":
            j += 1
        length = int(s[i:j])
        start = j + 1
        result.append(s[start : start + length])
        i = start + length
    return result
