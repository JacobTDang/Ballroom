#include <string>
#include <vector>

// EvalRPN evaluates an arithmetic expression given in Reverse Polish
// Notation and returns the result.
int EvalRPN(const std::vector<std::string>& tokens) {
    std::vector<int> stack;
    for (const auto& tok : tokens) {
        if (tok == "+" || tok == "-" || tok == "*" || tok == "/") {
            int b = stack.back();
            stack.pop_back();
            int a = stack.back();
            stack.pop_back();
            int res = 0;
            if (tok == "+") res = a + b;
            else if (tok == "-") res = a - b;
            else if (tok == "*") res = a * b;
            else res = a / b;
            stack.push_back(res);
        } else {
            stack.push_back(std::stoi(tok));
        }
    }
    return stack.back();
}
