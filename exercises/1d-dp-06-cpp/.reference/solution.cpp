#include <string>

// countExpansions grows outward from the center between indices l and
// r, counting one palindrome for every successful expansion.
static int countExpansions(const std::string& s, int l, int r) {
    int count = 0;
    while (l >= 0 && r < (int)s.size() && s[l] == s[r]) {
        count++;
        l--;
        r++;
    }
    return count;
}

// CountSubstrings returns the number of palindromic substrings in s,
// counting substrings at different positions separately even if they
// contain the same characters.
int CountSubstrings(std::string s) {
    int count = 0;
    for (int i = 0; i < (int)s.size(); i++) {
        count += countExpansions(s, i, i);
        count += countExpansions(s, i, i + 1);
    }
    return count;
}
