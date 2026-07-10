#include <functional>
#include <string>
#include <vector>

// Partition returns every way to split s into substrings that are
// all palindromes.
std::vector<std::vector<std::string>> Partition(std::string s) {
    auto isPalindrome = [](const std::string& str) {
        int l = 0, r = static_cast<int>(str.size()) - 1;
        while (l < r) {
            if (str[l] != str[r]) return false;
            l++;
            r--;
        }
        return true;
    };

    std::vector<std::vector<std::string>> res;
    std::vector<std::string> cur;
    std::function<void(int)> backtrack = [&](int start) {
        if (start == static_cast<int>(s.size())) {
            res.push_back(cur);
            return;
        }
        for (int end = start + 1; end <= static_cast<int>(s.size()); end++) {
            std::string sub = s.substr(start, end - start);
            if (isPalindrome(sub)) {
                cur.push_back(sub);
                backtrack(end);
                cur.pop_back();
            }
        }
    };
    backtrack(0);
    return res;
}
