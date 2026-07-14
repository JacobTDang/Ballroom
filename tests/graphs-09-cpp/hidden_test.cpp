#include <cassert>
#include <cstdio>
#include <map>
#include <set>
#include <vector>

std::vector<int> FindOrder(int numCourses, std::vector<std::vector<int>>& prerequisites);

bool isValidOrder(int numCourses, const std::vector<std::vector<int>>& prerequisites,
                   const std::vector<int>& order) {
    if (static_cast<int>(order.size()) != numCourses) return false;
    std::map<int, int> pos;
    std::set<int> seen;
    for (int i = 0; i < static_cast<int>(order.size()); i++) {
        if (seen.count(order[i])) return false;
        seen.insert(order[i]);
        pos[order[i]] = i;
    }
    for (auto& p : prerequisites) {
        int course = p[0], pre = p[1];
        if (pos[pre] >= pos[course]) return false;
    }
    return true;
}

void testValid() {
    std::vector<std::vector<int>> prereqs = {{1, 0}, {2, 0}, {3, 1}, {3, 2}};
    auto order = FindOrder(4, prereqs);
    assert(isValidOrder(4, prereqs, order));
}

void testCycle() {
    std::vector<std::vector<int>> prereqs = {{1, 0}, {0, 1}};
    auto order = FindOrder(2, prereqs);
    assert(order.empty());
}

void testNoPrerequisites() {
    std::vector<std::vector<int>> prereqs = {};
    auto order = FindOrder(3, prereqs);
    assert(isValidOrder(3, prereqs, order));
}

void testLinearChain() {
    std::vector<std::vector<int>> prereqs = {{1, 0}, {2, 1}, {3, 2}, {4, 3}};
    auto order = FindOrder(5, prereqs);
    assert(isValidOrder(5, prereqs, order));
}

int main() {
    testValid();
    testCycle();
    testNoPrerequisites();
    testLinearChain();
    std::printf("all tests passed\n");
    return 0;
}
