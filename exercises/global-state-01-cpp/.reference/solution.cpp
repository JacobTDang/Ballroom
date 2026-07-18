#include <string>
#include <vector>

// Formats items (low-stock item names) into report lines, one per
// item.
std::vector<std::string> generate_report(const std::vector<std::string>& items) {
    std::vector<std::string> report_lines;
    for (const auto& item : items) {
        report_lines.push_back("LOW STOCK: " + item);
    }
    return report_lines;
}
