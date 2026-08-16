# Reverse Words, Keep Punctuation Simple

A notification service flips message text word-by-word for an internal
readability experiment. Any punctuation stays attached to the word it
was typed with; only the order of the words changes.

Write `reverse_words(sentence: str) -> str` returning the words of the
sentence in reverse order, joined by single spaces. The input uses
exactly one space between words and has no leading or trailing
whitespace.

## Examples

    reverse_words("pay the balance now")  -> "now balance the pay"
    reverse_words("done!")                -> "done!"
    reverse_words("a b")                  -> "b a"

## Constraints

- 1 <= len(sentence) <= 1_000
- words are separated by exactly one space
