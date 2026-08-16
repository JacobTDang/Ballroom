def reverse_words(sentence: str) -> str:
    """Return the sentence with its word order reversed."""
    return " ".join(reversed(sentence.split(" ")))
