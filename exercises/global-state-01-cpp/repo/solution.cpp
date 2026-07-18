#include <string>
#include <vector>

// Formats items (low-stock item names) into report lines, one per
// item. Currently a call's report can contain lines left over from an
// earlier call -- find and fix the bug.
std::vector<std::string> generate_report(const std::vector<std::string>& items) {
    static std::vector<std::string> report_lines;
    for (const auto& item : items) {
        report_lines.push_back("LOW STOCK: " + item);
    }
    return report_lines;
}
