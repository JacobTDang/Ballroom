#include <map>
#include <sstream>
#include <string>

static std::string Trim(const std::string& s) {
    size_t a = s.find_first_not_of(" \t\r");
    if (a == std::string::npos) return "";
    size_t b = s.find_last_not_of(" \t\r");
    return s.substr(a, b - a + 1);
}

// Line-oriented with a current-section cursor. Every branch either
// consumes the line's whole meaning or errors with its 1-based number.
bool Parse(const std::string& input,
           std::map<std::string, std::map<std::string, std::string>>* out,
           std::string* err) {
    (*out)[""];
    std::string section = "";
    std::istringstream ss(input);
    std::string raw;
    int n = 0;
    while (std::getline(ss, raw)) {
        n++;
        std::string line = Trim(raw);
        if (line.empty() || line[0] == '#' || line[0] == ';') continue;
        if (line[0] == '[') {
            if (line.back() != ']') {
                *err = "line " + std::to_string(n) + ": unclosed section header";
                return false;
            }
            section = Trim(line.substr(1, line.size() - 2));
            (*out)[section];
        } else if (line.find('=') != std::string::npos) {
            auto eq = line.find('=');
            std::string key = Trim(line.substr(0, eq));
            if (key.empty()) {
                *err = "line " + std::to_string(n) + ": empty key";
                return false;
            }
            (*out)[section][key] = Trim(line.substr(eq + 1)); // later wins
        } else {
            *err = "line " + std::to_string(n) + ": not a header, comment, or key=value";
            return false;
        }
    }
    return true;
}
