#include <string>
#include <vector>

// Encodes a list of strings into a single string that decode can
// reconstruct exactly, including strings that contain any characters
// (digits, delimiters, etc). Each string is prefixed with its length and
// a '#' delimiter, so the delimiter itself can safely appear inside a
// string without ambiguity.
std::string encode(const std::vector<std::string>& strs) {
    std::string result;
    for (const auto& s : strs) {
        result += std::to_string(s.size());
        result += '#';
        result += s;
    }
    return result;
}

// Reverses encode.
std::vector<std::string> decode(const std::string& s) {
    std::vector<std::string> result;
    size_t i = 0;
    while (i < s.size()) {
        size_t j = i;
        while (s[j] != '#') j++;
        int length = std::stoi(s.substr(i, j - i));
        size_t start = j + 1;
        result.push_back(s.substr(start, length));
        i = start + length;
    }
    return result;
}
