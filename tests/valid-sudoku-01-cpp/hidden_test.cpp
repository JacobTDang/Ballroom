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
    std::vector<std::string> solved_board = {
        "534678912", "672195348", "198342567", "859761423", "426853791",
        "713924856", "961537284", "287419635", "345286179"};
    std::vector<std::string> same_digit_different_units_board = {
        "5........", ".........", ".........", ".........", "....5....",
        ".........", ".........", ".........", "........."};
    std::vector<std::string> single_cell_board = {
        "5........", ".........", ".........", ".........", ".........",
        ".........", ".........", ".........", "........."};

    assert(is_valid_sudoku(valid_board) == true);
    assert(is_valid_sudoku(invalid_column_board) == false);
    assert(is_valid_sudoku(invalid_row_board) == false);
    assert(is_valid_sudoku(invalid_box_board) == false);
    assert(is_valid_sudoku(empty_board) == true);
    assert(is_valid_sudoku(solved_board) == true);
    assert(is_valid_sudoku(same_digit_different_units_board) == true);
    assert(is_valid_sudoku(single_cell_board) == true);
    printf("all assertions passed\n");
    return 0;
}
