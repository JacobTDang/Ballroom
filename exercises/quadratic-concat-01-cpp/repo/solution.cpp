#include <string>
#include <vector>

// Joins log chunks that arrive newest-first into a single oldest-first
// log (no separator between chunks -- each chunk already carries its
// own formatting). Currently far slower than it should be on a large
// page -- find and fix the bug.
std::string build_log(const std::vector<std::string>& chunks) {
    std::string s;
    for (const auto& c : chunks) {
        s = c + s;
    }
    return s;
}
