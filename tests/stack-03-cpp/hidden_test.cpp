#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

int EvalRPN(const std::vector<std::string>& tokens);

int main() {
    assert(EvalRPN({"2", "1", "+", "3", "*"}) == 9);
    assert(EvalRPN({"4", "13", "5", "/", "+"}) == 6);
    assert(EvalRPN({"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"}) == 22);
    assert(EvalRPN({"18"}) == 18);
    assert(EvalRPN({"4", "3", "-"}) == 1);
    assert(EvalRPN({"-3", "4", "+"}) == 1);
    assert(EvalRPN({"7", "-2", "/"}) == -3);
    assert(EvalRPN({"5", "5", "*", "5", "*"}) == 125);
    printf("all assertions passed\n");
    return 0;
}
