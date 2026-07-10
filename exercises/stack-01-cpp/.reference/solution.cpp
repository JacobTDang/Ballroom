#include <string>
#include <unordered_map>
#include <vector>

// IsValid reports whether s's brackets are balanced and correctly
// nested.
bool IsValid(const std::string& s) {
    std::unordered_map<char, char> pairs = {{')', '('}, {']', '['}, {'}', '{'}};
    std::vector<char> stack;
    for (char c : s) {
        auto it = pairs.find(c);
        if (it != pairs.end()) {
            if (stack.empty() || stack.back() != it->second) return false;
            stack.pop_back();
        } else {
            stack.push_back(c);
        }
    }
    return stack.empty();
}
