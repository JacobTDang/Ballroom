#include <string>

// Whole-string glob: * (any run), ? (exactly one), [a-c] (one from a
// set/range). An unclosed [ makes the pattern match nothing.
//
// TODO: this handles only a lone "*" and literal equality.
bool Match(const std::string& pattern, const std::string& s) {
    if (pattern == "*") return true;
    return pattern == s;
}
