#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

bool is_valid_sudoku(const std::vector<std::string>& board);

int main() {
    std::vector<std::string> valid_board = {
        "53..7....", "6..195...", ".98....6.", "8...6...3", "4..8.3..1",
        "7...2...6", ".6....28.", "...419..5", "....8..79"};
    std::vector<std::string> invalid_column_board = {
        "83..7....", "6..195...", ".98....6.", "8...6...3", "4..8.3..1",
        "7...2...6", ".6....28.", "...419..5", "....8..79"};
    std::vector<std::string> invalid_row_board = {
        "5.......5", ".........", ".........", ".........", ".........",
        ".........", ".........", ".........", "........."};
    std::vector<std::string> invalid_box_board = {
        "1........", ".1.......", ".........", ".........", ".........",
        ".........", ".........", ".........", "........."};
    std::vector<std::string> empty_board = {
        ".........", ".........", ".........", ".........", ".........",
        ".........", ".........", ".........", "........."};

    assert(is_valid_sudoku(valid_board) == true);
    assert(is_valid_sudoku(invalid_column_board) == false);
    assert(is_valid_sudoku(invalid_row_board) == false);
    assert(is_valid_sudoku(invalid_box_board) == false);
    assert(is_valid_sudoku(empty_board) == true);
    printf("all assertions passed\n");
    return 0;
}
