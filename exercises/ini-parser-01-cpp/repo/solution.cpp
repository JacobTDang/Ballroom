#include <map>
#include <sstream>
#include <string>

// Parse reads an INI document into (*out)[section][key] = value.
//
// TODO: no sections, no comments, no errors -- every line is split on
// '=' into the "" section, and malformed lines are silently skipped.
bool Parse(const std::string& input,
           std::map<std::string, std::map<std::string, std::string>>* out,
           std::string* err) {
    std::istringstream ss(input);
    std::string line;
    while (std::getline(ss, line)) {
        auto eq = line.find('=');
        if (eq != std::string::npos) {
            (*out)[""][line.substr(0, eq)] = line.substr(eq + 1);
        }
    }
    return true;
}
