#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

int NumIslands(std::vector<std::vector<char>>& grid);

std::vector<std::vector<char>> gridOf(std::vector<std::string> rows) {
    std::vector<std::vector<char>> g;
    for (auto& r : rows) g.emplace_back(r.begin(), r.end());
    return g;
}

int main() {
    auto g1 = gridOf({"11110", "11010", "11000", "00000"});
    assert(NumIslands(g1) == 1);
    auto g2 = gridOf({"11000", "11000", "00100", "00011"});
    assert(NumIslands(g2) == 3);
    auto g3 = gridOf({"0"});
    assert(NumIslands(g3) == 0);
    auto g4 = gridOf({"1"});
    assert(NumIslands(g4) == 1);
    auto g5 = gridOf({"000", "000"});
    assert(NumIslands(g5) == 0);
    auto g6 = gridOf({"11", "11"});
    assert(NumIslands(g6) == 1);
    printf("all assertions passed\n");
    return 0;
}
