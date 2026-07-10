#include <cctype>
#include <string>

// IsPalindrome reports whether s is a palindrome, considering only
// alphanumeric characters and ignoring case.
bool IsPalindrome(const std::string& s) {
    int lo = 0, hi = static_cast<int>(s.size()) - 1;
    while (lo < hi) {
        while (lo < hi && !std::isalnum(static_cast<unsigned char>(s[lo]))) lo++;
        while (lo < hi && !std::isalnum(static_cast<unsigned char>(s[hi]))) hi--;
        if (std::tolower(static_cast<unsigned char>(s[lo])) !=
            std::tolower(static_cast<unsigned char>(s[hi]))) {
            return false;
        }
        lo++;
        hi--;
    }
    return true;
}
