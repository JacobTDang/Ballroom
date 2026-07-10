#include <functional>
#include <string>
#include <vector>

// GenerateParenthesis returns every well-formed combination of n pairs
// of parentheses.
std::vector<std::string> GenerateParenthesis(int n) {
    std::vector<std::string> res;
    std::function<void(std::string, int, int)> backtrack =
        [&](std::string cur, int open, int close) {
            if (static_cast<int>(cur.size()) == 2 * n) {
                res.push_back(cur);
                return;
            }
            if (open < n) backtrack(cur + "(", open + 1, close);
            if (close < open) backtrack(cur + ")", open, close + 1);
        };
    backtrack("", 0, 0);
    return res;
}
