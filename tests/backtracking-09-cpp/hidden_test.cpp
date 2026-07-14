#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::vector<std::string>> SolveNQueens(int n);

std::vector<std::vector<std::string>> normalizeExact(std::vector<std::vector<std::string>> boards) {
    std::sort(boards.begin(), boards.end());
    return boards;
}

int main() {
    {
        auto got = normalizeExact(SolveNQueens(4));
        auto want = normalizeExact({
            {".Q..", "...Q", "Q...", "..Q."},
            {"..Q.", "Q...", "...Q", ".Q.."},
        });
        assert(got == want);
    }
    {
        auto got = SolveNQueens(1);
        assert((got == std::vector<std::vector<std::string>>{{"Q"}}));
    }
    assert(SolveNQueens(2).empty());
    assert(SolveNQueens(3).empty());
    {
        auto got = SolveNQueens(5);
        assert(got.size() == 10);
        for (const auto& board : got) {
            assert(board.size() == 5);
            for (const auto& row : board) {
                assert(row.size() == 5);
            }
        }
    }
    printf("all assertions passed\n");
    return 0;
}
