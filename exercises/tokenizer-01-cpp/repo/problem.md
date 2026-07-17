# Arithmetic Tokenizer

Turn `"price * (1 + tax_rate)"` into a token stream a parser could
consume: numbers (with at most one decimal point), identifiers
(letter/underscore then letters/digits/underscores), the operators
`+ - * /`, and parentheses — each token carrying its **byte position**.
Whitespace separates; anything else is a loud error naming its
position.

The starter splits on spaces and calls everything an identifier.

## The invariant the tests enforce

- Exact token kinds, texts, and positions — including inputs with no
  spaces at all (`"3+4.5*x"`).
- `12..3` is an error (second dot), `@` is an error at its position,
  empty input is an empty token list.

API: `struct Token { std::string kind, text; int pos; }` and `std::vector<Token> Tokenize(const std::string& input, std::string* err)` — on error, fill `*err` (naming the position) and return an empty vector.
