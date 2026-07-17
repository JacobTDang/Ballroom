# INI Parser

Parse the classic config format: `[section]` headers, `key = value`
pairs, full-line comments starting with `#` or `;`, blank lines,
whitespace trimmed around keys and values. Keys before any section
header belong to the `""` section; a duplicate key's **later** value
wins.

Malformed input — a line with no `=` that isn't a header or comment,
or an unclosed `[section` — is an **error naming the line number**,
never silently skipped.

## The invariant the tests enforce

Exact parsed structure for a real config; each rule above has a case;
error cases name their 1-based line numbers.

API: `bool Parse(const std::string& input, std::map<std::string, std::map<std::string, std::string>>* out, std::string* err)` — on malformed input, fill `*err` naming the 1-based line number and return false.
