#include <cctype>
#include <string>
#include <vector>

struct Token {
    std::string kind;
    std::string text;
    int pos;
};

// One pass, one branch per character class -- each consumes a maximal
// token and records where it started. Errors name the exact position.
std::vector<Token> Tokenize(const std::string& input, std::string* err) {
    std::vector<Token> tokens;
    size_t i = 0;
    const size_t n = input.size();
    while (i < n) {
        char c = input[i];
        if (c == ' ' || c == '\t' || c == '\n') {
            i++;
        } else if (isdigit((unsigned char)c)) {
            size_t start = i;
            bool saw_dot = false;
            while (i < n && (isdigit((unsigned char)input[i]) || input[i] == '.')) {
                if (input[i] == '.') {
                    if (saw_dot) {
                        *err = "second decimal point at position " + std::to_string(i);
                        return {};
                    }
                    saw_dot = true;
                }
                i++;
            }
            tokens.push_back({"number", input.substr(start, i - start), (int)start});
        } else if (isalpha((unsigned char)c) || c == '_') {
            size_t start = i;
            while (i < n && (isalnum((unsigned char)input[i]) || input[i] == '_')) i++;
            tokens.push_back({"ident", input.substr(start, i - start), (int)start});
        } else if (c == '+' || c == '-' || c == '*' || c == '/') {
            tokens.push_back({"op", std::string(1, c), (int)i});
            i++;
        } else if (c == '(') {
            tokens.push_back({"lparen", "(", (int)i});
            i++;
        } else if (c == ')') {
            tokens.push_back({"rparen", ")", (int)i});
            i++;
        } else {
            *err = std::string("unexpected character '") + c + "' at position " + std::to_string(i);
            return {};
        }
    }
    return tokens;
}
