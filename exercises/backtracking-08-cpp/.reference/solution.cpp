#include <functional>
#include <string>
#include <unordered_map>
#include <vector>

// LetterCombinations returns every letter combination that digits
// could represent on a phone keypad.
std::vector<std::string> LetterCombinations(std::string digits) {
    if (digits.empty()) return {};
    static const std::unordered_map<char, std::string> phoneLetters = {
        {'2', "abc"}, {'3', "def"}, {'4', "ghi"}, {'5', "jkl"},
        {'6', "mno"}, {'7', "pqrs"}, {'8', "tuv"}, {'9', "wxyz"},
    };

    std::vector<std::string> res;
    std::string cur;
    std::function<void(int)> backtrack = [&](int idx) {
        if (idx == static_cast<int>(digits.size())) {
            res.push_back(cur);
            return;
        }
        const std::string& letters = phoneLetters.at(digits[idx]);
        for (char c : letters) {
            cur.push_back(c);
            backtrack(idx + 1);
            cur.pop_back();
        }
    };
    backtrack(0);
    return res;
}
