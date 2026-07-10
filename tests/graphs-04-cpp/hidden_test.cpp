#include <cassert>
#include <cstdio>
#include <vector>

const int kInf = 2147483647;

void WallsAndGates(std::vector<std::vector<int>>& rooms);

int main() {
    {
        std::vector<std::vector<int>> rooms = {
            {kInf, -1, 0, kInf},
            {kInf, kInf, kInf, -1},
            {kInf, -1, kInf, -1},
            {0, -1, kInf, kInf},
        };
        std::vector<std::vector<int>> want = {
            {3, -1, 0, 1},
            {2, 2, 1, -1},
            {1, -1, 2, -1},
            {0, -1, 3, 4},
        };
        WallsAndGates(rooms);
        assert(rooms == want);
    }
    {
        std::vector<std::vector<int>> rooms = {{0, -1, kInf}};
        std::vector<std::vector<int>> want = {{0, -1, kInf}};
        WallsAndGates(rooms);
        assert(rooms == want);
    }
    {
        std::vector<std::vector<int>> rooms = {{kInf, kInf}};
        std::vector<std::vector<int>> want = {{kInf, kInf}};
        WallsAndGates(rooms);
        assert(rooms == want);
    }
    printf("all assertions passed\n");
    return 0;
}
