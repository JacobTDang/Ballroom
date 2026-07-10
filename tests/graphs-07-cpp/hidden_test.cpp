#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

void Solve(std::vector<std::vector<char>>& board);

std::vector<std::vector<char>> toGrid(const std::vector<std::string>& rows) {
    std::vector<std::vector<char>> out;
    for (const auto& row : rows) {
        out.push_back(std::vector<char>(row.begin(), row.end()));
    }
    return out;
}

void testClassic() {
    auto board = toGrid({"XXXX", "XOOX", "XXOX", "XOXX"});
    auto want = toGrid({"XXXX", "XXXX", "XXXX", "XOXX"});
    Solve(board);
    assert(board == want);
}

void testAllBorderConnected() {
    auto board = toGrid({"OOO", "OXO", "OOO"});
    auto want = toGrid({"OOO", "OXO", "OOO"});
    Solve(board);
    assert(board == want);
}

void testSingleCell() {
    auto board = toGrid({"O"});
    auto want = toGrid({"O"});
    Solve(board);
    assert(board == want);
}

int main() {
    testClassic();
    testAllBorderConnected();
    testSingleCell();
    std::printf("all tests passed\n");
    return 0;
}
