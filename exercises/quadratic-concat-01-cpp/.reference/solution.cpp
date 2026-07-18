#include <string>
#include <vector>

// Joins log chunks that arrive newest-first into a single oldest-first
// log (no separator between chunks -- each chunk already carries its
// own formatting).
std::string build_log(const std::vector<std::string>& chunks) {
    std::string s;
    for (auto it = chunks.rbegin(); it != chunks.rend(); ++it) {
        s += *it;
    }
    return s;
}
