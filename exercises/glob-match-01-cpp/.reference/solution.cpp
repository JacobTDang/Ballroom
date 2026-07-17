#include <string>

// Match c against the class starting at pattern[p] == '['. Sets
// matched and next (index past ']'); returns false for an unclosed
// class.
static bool MatchClass(const std::string& pattern, size_t p, char c,
                       bool* matched, size_t* next) {
    size_t q = p + 1;
    *matched = false;
    const size_t n = pattern.size();
    while (q < n && pattern[q] != ']') {
        if (q + 2 < n && pattern[q + 1] == '-' && pattern[q + 2] != ']') {
            if (pattern[q] <= c && c <= pattern[q + 2]) *matched = true;
            q += 3;
        } else {
            if (pattern[q] == c) *matched = true;
            q++;
        }
    }
    if (q >= n) return false; // unclosed
    *next = q + 1;
    return true;
}

// The classic two-pointer loop: on '*' remember both positions; on a
// later mismatch, back up to just after the star and let it swallow
// one more character. That remembered pair IS the backtracking state.
bool Match(const std::string& pattern, const std::string& s) {
    size_t p = 0, i = 0;
    size_t star_p = std::string::npos, star_i = 0;

    while (i < s.size()) {
        bool advanced = false;
        if (p < pattern.size()) {
            char ch = pattern[p];
            if (ch == '*') {
                star_p = p;
                star_i = i;
                p++;
                continue;
            }
            // The class branch must run before the literal branch: a
            // pattern "[" against the string "[" would otherwise match
            // itself literally instead of being an invalid class.
            if (ch == '[') {
                bool matched;
                size_t next;
                if (!MatchClass(pattern, p, s[i], &matched, &next)) return false;
                if (matched) {
                    p = next;
                    i++;
                    advanced = true;
                }
            } else if (ch == '?' || ch == s[i]) {
                p++;
                i++;
                advanced = true;
            }
        }
        if (advanced) continue;
        if (star_p == std::string::npos) return false;
        star_i++;
        p = star_p + 1;
        i = star_i;
    }

    while (p < pattern.size() && pattern[p] == '*') p++;
    return p == pattern.size();
}
