#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

bool Exist(std::vector<std::vector<char>>& board, std::string word);

std::vector<std::vector<char>> makeBoard() {
    return {{'A', 'B', 'C', 'E'}, {'S', 'F', 'C', 'S'}, {'A', 'D', 'E', 'E'}};
}

int main() {
    {
        auto b = makeBoard();
        assert(Exist(b, "ABCCED") == true);
    }
    {
        auto b = makeBoard();
        assert(Exist(b, "SEE") == true);
    }
    {
        auto b = makeBoard();
        assert(Exist(b, "ABCB") == false);
    }
    {
        auto b = makeBoard();
        assert(Exist(b, "ABFSAB") == false);
    }
    {
        std::vector<std::vector<char>> b = {{'A'}};
        assert(Exist(b, "A") == true);
    }
    {
        std::vector<std::vector<char>> b = {{'A'}};
        assert(Exist(b, "AA") == false);
    }
    {
        std::vector<std::vector<char>> b = {{'a', 'b'}, {'c', 'd'}};
        assert(Exist(b, "abdc") == true);
    }
    {
        std::vector<std::vector<char>> b = {{'a', 'b'}, {'c', 'd'}};
        assert(Exist(b, "abcd") == false);
    }
    printf("all assertions passed\n");
    return 0;
}
