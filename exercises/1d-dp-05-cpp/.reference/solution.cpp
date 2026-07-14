#include <string>

// expandAroundCenter grows outward from the center between indices l
// and r (l == r for an odd-length center, r == l+1 for an even-length
// center) and returns the length of the palindrome found.
static int expandAroundCenter(const std::string& s, int l, int r) {
    while (l >= 0 && r < (int)s.size() && s[l] == s[r]) {
        l--;
        r++;
    }
    return r - l - 1;
}

// LongestPalindrome returns the longest palindromic substring of s. If
// several substrings share the maximum length, the first one found
// scanning left to right is returned.
std::string LongestPalindrome(std::string s) {
    if (s.empty()) return "";
    int start = 0, end = 0;
    for (int i = 0; i < (int)s.size(); i++) {
        int len1 = expandAroundCenter(s, i, i);
        int len2 = expandAroundCenter(s, i, i + 1);
        int maxLen = len1 > len2 ? len1 : len2;
        if (maxLen > end - start + 1) {
            start = i - (maxLen - 1) / 2;
            end = i + maxLen / 2;
        }
    }
    return s.substr(start, end - start + 1);
}
